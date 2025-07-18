import logging
import random
import time

from selenium.common import TimeoutException, NoSuchElementException


from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

from browser import BrowserManager


class TelegramBotAutomation:
    def __init__(self, serial_number):
        self.serial_number = serial_number
        self.browser_manager = BrowserManager(serial_number)
        logging.info(f"Initializing automation for account {serial_number}")
        self.browser_manager.start_browser()
        self.driver = self.browser_manager.driver

    def navigate_to_bot(self):
        try:
            self.driver.get('https://web.telegram.org/a/')
            logging.info(f"Account {self.serial_number}: Navigated to Telegram web.")
        except Exception as e:
            logging.exception(f"Account {self.serial_number}: Exception in navigating to Telegram bot: {str(e)}")
            self.browser_manager.close_browser()

    def send_message(self, message):
        chat_input_area = self.wait_for_element(By.XPATH,
                                                '/html/body/div[2]/div/div[1]/div/div[1]/div/div[2]/input')
        chat_input_area.click()
        chat_input_area.send_keys(message)

        search_area = self.wait_for_element(By.XPATH,
                                            '/html[1]/body[1]/div[1]/div[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[2]/div[2]/div[1]/div[1]/div[1]/div[2]/ul[1]/a[1]/div[1]')
        search_area.click()
        logging.info(f"Account {self.serial_number}: Group searched.")

    def click_link(self):
        link = self.wait_for_element(By.CSS_SELECTOR, "a[href*='t.me/Tomarket_ai_bot/app?startapp=0001Fpms']")
        link.click()

        # После последней обновы 31.08.24 закомментил
        # launch_click = self.wait_for_element(By.XPATH, "//body/div[@class='popup popup-peer popup-confirmation active']/div[@class='popup-container z-depth-1']/div[@class='popup-buttons']/button[1]/div[1]")
        # launch_click.click()
        logging.info(f"Account {self.serial_number}: TOMARKET STARTED")
        sleep_time = random.randint(20, 30)
        logging.info(f"Sleeping for {sleep_time} seconds.")
        time.sleep(sleep_time)
        if not self.switch_to_iframe():
            logging.info(f"Account {self.serial_number}: No iframes found")
            return

        try:
            daily_reward_button = WebDriverWait(self.driver, 1).until(
                EC.element_to_be_clickable(
                    (By.XPATH, "/html/body/div[2]/div/div/div[2]/div[2]/div[5]/div"))
            )
            daily_reward_button.click()
            logging.info(f"Account {self.serial_number}: Daily reward claimed.")
            time.sleep(2)
        except TimeoutException:
            logging.info(f"Account {self.serial_number}: Daily reward has already been claimed or button not found.")

    def check_claim_button(self):
        if not self.switch_to_iframe():
            logging.info(f"Account {self.serial_number}: No iframes found")
            return 0.0

        initial_balance = self.check_balance()
        self.process_buttons()
        final_balance = self.check_balance()

        if final_balance is not None and initial_balance == final_balance and not self.is_farming_active():
            raise Exception(f"Account {self.serial_number}: Balance did not change after claiming tokens.")

        return final_balance if final_balance is not None else 0.0

    def switch_to_iframe(self):
        self.driver.switch_to.default_content()
        iframes = self.driver.find_elements(By.TAG_NAME, "iframe")
        if iframes:
            self.driver.switch_to.frame(iframes[0])
            return True
        return False

    def process_buttons(self):
        parent_selector = "div.kit-fixed-wrapper.has-layout-tabs"

        button_primary_selector = "button.kit-button.is-large.is-primary.is-fill.button"
        button_done_selector = "button.kit-button.is-large.is-drop.is-fill.button.is-done"
        button_secondary_selector = "button.kit-button.is-large.is-secondary.is-fill.is-centered.button.is-active"

        parent_element = self.wait_for_element(By.CSS_SELECTOR, parent_selector)

        if parent_element:
            primary_buttons = parent_element.find_elements(By.CSS_SELECTOR, button_primary_selector)
            done_buttons = parent_element.find_elements(By.CSS_SELECTOR, button_done_selector)
            secondary_buttons = parent_element.find_elements(By.CSS_SELECTOR, button_secondary_selector)

            # logging.info(f"Account {self.serial_number}: Found {len(primary_buttons)} default button")
            # logging.info(f"Account {self.serial_number}: Found {len(done_buttons)} done button")
            # logging.info(f"Account {self.serial_number}: Found {len(secondary_buttons)} active button")

            for button in primary_buttons:
                self.process_single_button(button)
            for button in done_buttons:
                self.process_single_button(button)
            for button in secondary_buttons:
                self.process_single_button(button)
        else:
            logging.info(f"Account {self.serial_number}: Parent element not found.")

    def process_single_button(self, button):
        button_text = self.get_button_text(button)
        amount_elements = button.find_elements(By.CSS_SELECTOR, "div.amount")
        amount_text = amount_elements[0].text if amount_elements else None

        if "Farming" in button_text:
            self.handle_farming(button)
        elif "Start farming" in button_text and not amount_text:
            self.start_farming(button)
        elif amount_text:
            self.claim_tokens(button, amount_text)

    def get_button_text(self, button):
        try:
            return button.find_element(By.CSS_SELECTOR, ".button-label").text
        except NoSuchElementException:
            return button.find_element(By.CSS_SELECTOR, ".label").text

    def handle_farming(self, button):
        logging.info(
            f"Account {self.serial_number}: Farming is active. The account is currently farming. Checking timer again.")
        try:
            time_left = self.driver.find_element(By.CSS_SELECTOR, "div.time-left").text
            logging.info(f"Account {self.serial_number}: Remaining time to next claim opportunity: {time_left}")
        except NoSuchElementException:
            logging.warning(f"Account {self.serial_number}: Timer not found after detecting farming status.")

    def start_farming(self, button):
        button.click()
        logging.info(f"Account {self.serial_number}: Clicked on 'Start farming'.")
        sleep_time = random.randint(20, 30)
        logging.info(f"Sleeping for {sleep_time} seconds.")
        time.sleep(sleep_time)
        self.handle_farming(button)
        if not self.is_farming_active():
            raise Exception(f"Account {self.serial_number}: Farming did not start successfully.")

    def claim_tokens(self, button, amount_text):
        sleep_time = random.randint(5, 15)
        logging.info(f"Sleeping for {sleep_time} seconds.")
        time.sleep(sleep_time)
        logging.info(f"Account {self.serial_number}: Account has {amount_text} claimable tokens. Trying to claim.")

        button.click()
        logging.info(
            f"Account {self.serial_number}: Click successful. 10s sleep, waiting for button to update to 'Start Farming'...")
        time.sleep(10)

        WebDriverWait(self.driver, 10).until(
            EC.visibility_of_element_located((By.CSS_SELECTOR, ".label"))
        )

        start_farming_button = self.wait_for_element(By.CSS_SELECTOR, ".label")
        start_farming_button.click()
        logging.info(f"Account {self.serial_number}: Second click successful on 'Start farming'")
        sleep_time = random.randint(5, 15)
        logging.info(f"Sleeping for {sleep_time} seconds.")
        time.sleep(sleep_time)
        self.process_buttons()
        self.handle_farming(start_farming_button)
        if not self.is_farming_active():
            raise Exception(f"Account {self.serial_number}: Farming did not start successfully.")

    def check_balance(self):
        logging.info(f"Account {self.serial_number}: Trying to get total balance")
        try:
            iframes = self.driver.find_elements(By.TAG_NAME, "iframe")
            if iframes:
                self.driver.switch_to.frame(iframes[0])

            balance_elements = WebDriverWait(self.driver, 10).until(
                EC.visibility_of_all_elements_located((By.CSS_SELECTOR,
                                                       "div.profile-with-balance .kit-counter-animation.value .el-char-wrapper .el-char"))
            )
            balance = ''.join([element.text for element in balance_elements])
            logging.info(f"Account {self.serial_number}: Current balance: {balance}")
            sleep_time = random.randint(5, 15)
            logging.info(f"Sleeping for {sleep_time} seconds.")
            time.sleep(sleep_time)
            return float(balance.replace(',', ''))

        except TimeoutException:
            logging.warning(f"Account {self.serial_number}: Failed to find the balance element.")
            return 0.0

    def wait_for_element(self, by, value, timeout=10):
        return WebDriverWait(self.driver, timeout).until(
            EC.element_to_be_clickable((by, value))
        )

    def wait_for_elements(self, by, value, timeout=10):
        return WebDriverWait(self.driver, timeout).until(
            EC.visibility_of_all_elements_located((by, value))
        )

    def is_farming_active(self):
        try:
            self.driver.find_element(By.CSS_SELECTOR, "div.time-left")
            return True
        except NoSuchElementException:
            return False
