from selenium import webdriver
from selenium.webdriver.common.by import By
import requests
import json
import os


uuid = os.getenv("UUID")
if uuid is None:
    print("UUID is not set")
    print("Usage: UUID={uuid} python3 scraper.py")
    exit(1)
else:
    url = "https://schedule.mstimetables.ru/publications/" + uuid
    api_url = "https://schedule.mstimetables.ru/api/publications/" + uuid + "/"

if os.path.exists("groups.csv"):
    os.remove("groups.csv")
else:
    print("Groups.csv not found")
output_file = open("groups.csv", "w")
output_file.write("api_id,kind,name\n")

driver = webdriver.Firefox()

driver.get(url + "/#/groups")

groups = driver.find_elements(By.CLASS_NAME, "link")

for group in groups:
    name = group.text
    url = group.get_attribute("href")
    id = url.split("/")[7]
    print(id + ",group," + name)
    output_file.write(id + ",group," + name + "\n")

output_file.close()
driver.close()

if os.path.exists("teachers.csv"):
    print("found teachers.csv, removing")
    os.remove("teachers.csv")
else:
    print("teachers.csv not found")

output_file = open("teachers.csv", "w")
output_file.write("api_id,kind,name\n")

resp = requests.get(api_url + "teachers")
teachers = json.loads(resp.text)
for teacher in teachers:
    name = teacher["fio"]
    # name = teacher["fio"].split(" ")
    # while len(name) < 3:
    #     name.append("")
    id = str(teacher["id"])

    print(id + ",teacher," + name)
    output_file.write(id + ",teacher," + name + "\n")

output_file.close()
