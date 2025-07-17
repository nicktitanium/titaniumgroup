# import requests
# from selenium import webdriver
# from selenium.webdriver.common.by import By
# from selenium.webdriver.chrome.service import Service
#
# # Конфигурация API AdsPower
# ADSP_API_BASE = "http://local.adspower.net:50325"
# API_KEY = "6deaeb14ad322e86b2f45251b5d25f4d"
# PROFILE_ID = "23"
#
#
# def start_profile(profile_id):
#     # Запрос на запуск профиля
#     response = requests.get(f"{ADSP_API_BASE}/api/v1/browser/open?user_id={profile_id}&api_key={API_KEY}")
#     data = response.json()
#
#     if data["code"] == 0:
#         # Получаем URL для подключения к удаленному браузеру
#         remote_url = data["data"]["url"]
#         return remote_url
#     else:
#         print("Ошибка при запуске профиля:", data["msg"])
#         return None
#
#
# def open_youtube(remote_url):
#     # Подключаемся к удаленному браузеру через WebDriver
#     chrome_options = webdriver.ChromeOptions()
#     chrome_options.debugger_address = remote_url
#     driver = webdriver.Chrome(service=Service("/path/to/chromedriver"), options=chrome_options)
#
#     # Переходим на YouTube
#     driver.get("https://www.youtube.com")
#
#     # Дополнительные действия на YouTube
#     # Например, поиск по видео:
#     search_box = driver.find_element(By.NAME, "search_query")
#     search_box.send_keys("интересное видео")
#     search_box.submit()
#
#     # Закрыть браузер после выполнения действий
#     driver.quit()
#
#
# if __name__ == "__main__":
#     remote_url = start_profile(PROFILE_ID)
#     if remote_url:
#         open_youtube(remote_url)
import time

import requests
from selenium import webdriver
from selenium.webdriver import ActionChains
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service as ChromeService

# Настройки API AdsPower
API_BASE_URL = "http://local.adspower.net:50325"  # Локальный API AdsPower
PROFILE_ID = "km1gfn2"  # Уникальный ID профиля, замените на ваш
API_KEY = "78dc3887cd5e3293a89ad2e2d070f946"  # Замените на ваш API ключ

# Заголовки с API ключом
HEADERS = {
    "Authorization": f"Bearer {API_KEY}"
}


def check_api_status():
    response = requests.get(f"{API_BASE_URL}/status", headers=HEADERS)
    if response.json().get("code") == 0:
        print("API доступен.")
    else:
        print("API недоступен.")
    return response.json()


def start_browser(profile_id):
    response = requests.get(f"{API_BASE_URL}/api/v1/browser/start?user_id={profile_id}", headers=HEADERS)
    data = response.json()

    if data["code"] == 0:
        ws_url = data["data"]["ws"]["selenium"]  # URL WebSocket для автоматизации через Selenium
        webdriver_path = data["data"]["webdriver"]  # Путь к WebDriver
        print(f"Браузер запущен успешно. WebSocket URL: {ws_url}")
        return ws_url, webdriver_path
    else:
        print(f"Ошибка при запуске браузера: {data['msg']}")
        return None, None


def open_youtube(ws_url, webdriver_path):
    chrome_service = ChromeService(executable_path=webdriver_path)

    chrome_options = webdriver.ChromeOptions()
    chrome_options.debugger_address = ws_url

    driver = webdriver.Chrome(service=chrome_service, options=chrome_options)

    driver.get("https://web.telegram.org/a/#7200102626")

    # Найдите элемент кнопки по CSS_SELECTOR и нажмите на нее
    search_box = driver.find_element(By.CSS_SELECTOR,
                                     "#MiddleColumn > div.messages-layout > div.Transition > div > "
                                     "div.middle-column-footer > div > div > div > "
                                     "button.Button.bot-menu.open.default.translucent.round")
    # Прокрутка к кнопке, если она не видна
    actions = ActionChains(driver)
    actions.move_to_element(search_box).perform()

    # Нажмите на кнопку
    search_box.click()
    print("Кнопка успешно нажата.")

    # Закрытие текущей вкладки
    driver.close()

    # Завершение работы и закрытие браузера
    input("Нажмите Enter для завершения и закрытия браузера...")
    driver.quit()


def close_browser(profile_id):
    """Закрывает браузер в указанном профиле."""
    response = requests.get(f"{API_BASE_URL}/api/v1/browser/stop?user_id={profile_id}", headers=HEADERS)
    data = response.json()
    if data["code"] == 0:
        print("Браузер успешно закрыт.")
    else:
        print(f"Ошибка при закрытии браузера: {data['msg']}")


if __name__ == "__main__":

    check_api_status()

    ws_url, webdriver_path = start_browser(PROFILE_ID)
    if ws_url and webdriver_path:
        open_youtube(ws_url, webdriver_path)
        close_browser(PROFILE_ID)
