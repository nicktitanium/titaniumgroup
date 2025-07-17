import time

import base58
import requests
from environs import Env
from tronpy import Tron
from tronpy.keys import PrivateKey

env = Env()
env.read_env()

EXCEED_BALANCE = env.int('EXCEED_BALANCE')
FROM_ADDRESS = env('FROM_ADDRESS')
TO_ADDRESS = env('TO_ADDRESS')
FEE_LIMIT = env.int('FEE_LIMIT')
PRIVATE_KEY1 = env('PRIVATE_KEY1')
PRIVATE_KEY2 = env('PRIVATE_KEY2')
PERCENT_TO_SEND=env.float('PERCENT_TO_SEND')

def get_balance(address: str) -> int | None:
    try:
        url = "https://api.trongrid.io/wallet/getaccount"

        payload = {
            "address": base58.b58decode_check(address).hex(),
        }
        headers = {
            "accept": "application/json",
            "content-type": "application/json",
            'Tron-Pro-Api-Key': 'af5047d9-a3d0-4ca0-97b1-271dfe5da414'
        }

        response = requests.post(url, json=payload, headers=headers)
        data = response.json()
        if data:
            return int(data.get("balance"))
    except Exception as e:
        print("Ошибка: возвращаемые данные не содержат баланса")


def transfer(
        from_address: str,
        to_address: str,
        fee_limit: int,
        private_key1: str,
        private_key2: str,
        amount_to_send: int
):
    try:
        client = Tron(network='mainnet')
        private_key1_converted = PrivateKey(bytes.fromhex(private_key1))
        private_key2_converted = PrivateKey(bytes.fromhex(private_key2))

        txn = (
            client.trx.transfer(from_address, to_address, amount_to_send)
            .fee_limit(fee_limit)
            .memo("test memo")

        )

        txn_build = txn.build()
        txn_inspect = txn_build.inspect()
        txn_sign= txn_inspect.sign(private_key1_converted)
        txn_sign=txn_sign.sign(private_key2_converted)
        txn_broadcast= txn_sign.broadcast()

        print(txn_broadcast)
        print(txn_broadcast.wait())
    except Exception as e:
        print("Ошибка во время перевода:", e)
def main(from_address: str,
         to_address: str,
         exceed_balance: int,
         fee_limit: int,
         private_key1: str,
         private_key2: str,
         percent_to_send:float):
    print(f'Стартуем с  {from_address[:5]}... на адрес {to_address[:5]}..., Fee limit {fee_limit}, percent to send {percent_to_send}')
    while True:
        time.sleep(5)
        balance_got = get_balance(
            address=from_address
        )
        if not balance_got:
            print('Баланс не пришел')
            continue
        if not isinstance(balance_got, int):
            print('Баланс не INTEGER')
            continue
        print(f"Получен баланс {balance_got}")
        if balance_got > exceed_balance:
            amount_to_send = int(balance_got * percent_to_send)
            print(f"Полученный баланс выше диапазона, отправляется {amount_to_send}")
            transfer(
                from_address=from_address,
                to_address=to_address,
                fee_limit=fee_limit,
                private_key1=private_key1,
                private_key2=private_key2,
                amount_to_send=amount_to_send
            )


if __name__ == '__main__':
    main(
        from_address=FROM_ADDRESS,
        to_address=TO_ADDRESS,
        exceed_balance=EXCEED_BALANCE,
        fee_limit=FEE_LIMIT,
        private_key1=PRIVATE_KEY1,
        private_key2=PRIVATE_KEY2,
        percent_to_send=PERCENT_TO_SEND
    )
