import urllib.request
import urllib.error
import json
import time
import random

BASE_URL = "http://localhost:8088"

def make_request(method, endpoint, data=None, token=None):
    url = f"{BASE_URL}{endpoint}"
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    
    req_data = None
    if data:
        req_data = json.dumps(data).encode('utf-8')
    
    req = urllib.request.Request(url, data=req_data, headers=headers, method=method)
    
    try:
        with urllib.request.urlopen(req) as response:
            status = response.status
            resp_body = response.read().decode('utf-8')
            try:
                return status, json.loads(resp_body)
            except:
                return status, resp_body
    except urllib.error.HTTPError as e:
        resp_body = e.read().decode('utf-8')
        try:
            return e.code, json.loads(resp_body)
        except:
            return e.code, resp_body
    except Exception as e:
        return 0, str(e)

def run_tests():
    print("=== 开始全系统接口测试 ===")
    
    # 1. 注册普通用户
    ts = int(time.time())
    rnd = random.randint(1000, 9999)
    user_student_id = f"STU_{ts}_{rnd}"
    user_phone = f"138{ts}"[:11] # Phone might still collide if truncated? 
    # Let's make phone more unique too.
    user_phone = f"138{str(ts)[-4:]}{rnd}"
    
    user_pass = "password123"
    
    print(f"\n[1] 注册普通用户: {user_student_id}")
    reg_data = {
        "user_name": f"User_{ts}_{rnd}",
        "password": user_pass,
        "student_id": user_student_id,
        "phone": user_phone,
        "address": "Test Address",
        "role": "user"
    }
    status, resp = make_request("POST", "/register", reg_data)
    print(f"Status: {status}, Response: {resp}")
    if status != 200:
        print("注册失败，停止测试")
        return
    
    user_id = resp['user']['user_id']
    user_token = resp['access_token']
    print(f"User ID: {user_id}")
    
    # 2. 注册管理员用户
    admin_student_id = f"ADM_{ts}_{rnd}"
    admin_phone = f"139{str(ts)[-4:]}{rnd}"
    
    print(f"\n[2] 注册管理员用户: {admin_student_id}")
    admin_reg_data = {
        "user_name": f"Admin_{ts}_{rnd}",
        "password": user_pass,
        "student_id": admin_student_id,
        "phone": admin_phone,
        "address": "Admin Office",
        "role": "admin"
    }
    status, resp = make_request("POST", "/register", admin_reg_data)
    print(f"Status: {status}")
    if status != 200:
        print("管理员注册失败")
        return
    admin_token = resp['access_token']

    # 3. 用户登录 (验证登录接口)
    print(f"\n[3] 用户登录: {user_student_id}")
    login_data = {
        "student_id": user_student_id,
        "password": user_pass
    }
    status, resp = make_request("POST", "/login", login_data)
    print(f"Status: {status}")
    if status == 200:
        user_token = resp['access_token'] # Refresh token just in case
    
    # 4. 用户更新信息
    print(f"\n[4] 用户更新信息")
    update_info_data = {
        "user_id": user_id,
        "user_name": f"User_{ts}_Updated",
        "address": "New Address"
    }
    status, resp = make_request("POST", "/updateUserInfo", update_info_data, user_token)
    print(f"Status: {status}, Response: {resp}")

    # 5. 用户寄件 (MailPack)
    print(f"\n[5] 用户寄件")
    mail_data = {
        "shipping_address": "Sender Addr",
        "recipient": "Receiver Name",
        "reciving_address": "Receiver Addr",
        "shipper_phone": "13800000000",
        "recipient_phone": user_phone
    }
    status, resp = make_request("POST", "/mailPack", mail_data, user_token)
    print(f"Status: {status}, Response: {resp}")
    mail_pack_id = 0
    if status == 200:
        mail_pack_id = resp['pack']['pack_id']
        print(f"Mailed Pack ID: {mail_pack_id}")

    # 6. 用户取消寄件
    if mail_pack_id:
        print(f"\n[6] 用户取消寄件: {mail_pack_id}")
        cancel_data = {
            "pack_id": mail_pack_id,
            "user_id": user_id
        }
        status, resp = make_request("POST", "/cancelMail", cancel_data, user_token)
        print(f"Status: {status}, Response: {resp}")

    # 7. 包裹入库 (CheckInPack) - 模拟站点操作
    print(f"\n[7] 包裹入库 (模拟)")
    checkin_data = {
        "pack_id": 0, # Will be ignored/generated if 0? No, checkInPack usually takes an ID or generates one? 
                      # Looking at controller: CheckInPack takes PackId. 
                      # Usually CheckIn is for incoming packs. Let's use a random ID.
        "user_id": user_id,
        "shelf_code": 101
    }
    # Note: PackId in CheckInPack is usually the tracking number from external courier.
    # Since we don't have external tracking numbers, we'll use a random large int.
    external_pack_id = int(time.time() * 1000) 
    checkin_data['pack_id'] = external_pack_id
    
    status, resp = make_request("POST", "/packCheckIn", checkin_data, user_token)
    print(f"Status: {status}, Response: {resp}")
    
    # 8. 获取用户所有包裹
    print(f"\n[8] 获取用户所有包裹: {user_id}")
    status, resp = make_request("GET", f"/allPacks/{user_id}", None, user_token)
    print(f"Status: {status}")
    if status == 200:
        packs = resp.get('packs', [])
        print(f"Found {len(packs)} packs")

    # 9. 管理员获取所有用户
    print(f"\n[9] 管理员获取所有用户")
    status, resp = make_request("GET", "/admin/users", None, admin_token)
    print(f"Status: {status}")
    if status == 200:
        users = resp.get('users', [])
        print(f"Found {len(users)} users")

    # 10. 管理员获取所有包裹
    print(f"\n[10] 管理员获取所有包裹")
    status, resp = make_request("GET", "/admin/packs", None, admin_token)
    print(f"Status: {status}")
    if status == 200:
        all_packs = resp.get('packs', [])
        print(f"Found {len(all_packs)} packs in system")

    # 11. 管理员更新包裹信息
    print(f"\n[11] 管理员更新包裹信息: {external_pack_id}")
    update_pack_data = {
        "pack_id": external_pack_id,
        "pack_status": "arrived",
        "pickup_code": "999-8888"
    }
    status, resp = make_request("PUT", "/admin/pack", update_pack_data, admin_token)
    print(f"Status: {status}, Response: {resp}")

    # 12. 验证更新结果 (获取详情)
    print(f"\n[12] 验证更新结果")
    status, resp = make_request("GET", f"/getPackDetails/{external_pack_id}", None, user_token)
    print(f"Status: {status}, Response: {resp}")

    print("\n=== 测试结束 ===")

if __name__ == "__main__":
    run_tests()
