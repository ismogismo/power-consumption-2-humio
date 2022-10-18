# https://api.eloverblik.dk/CustomerApi/swagger/index.html
import requests
import json

# Read refresh token from file
f = open("rt.secret", "r")
token = f.readline()

# Get data access token for subsequent requests
get_data_access_token_url = 'https://api.eloverblik.dk/CustomerApi/api/token'
headers = {
    'accept': 'application/json',
    'Authorization': 'Bearer ' + token,
}

response = requests.get(get_data_access_token_url, headers=headers)
data_access_token = response.json()['result']

print('---- retrieved access token ----\n')

# Get id of first meter - edit if you have more than one meter
metering_points_url = 'https://api.eloverblik.dk/CustomerApi/api/meteringpoints/meteringpoints'
headers = {
    'accept': 'application/json',
    'Authorization': 'Bearer ' + data_access_token,
}
meters = requests.get(metering_points_url, headers=headers)
if meters.status_code != 200:
    print('FAILED getting metering id')
    print(meters.status_code)
    exit()

first_meter = meters.json()['result'][0]['meteringPointId']

print('---- retrieved meter id ----\n')
print('meter id: ')
print(first_meter)

#Try to get data
meter_data = 'https://api.eloverblik.dk/CustomerApi/api/meterdata/gettimeseries/'
timeseries_data = {
    'dateFrom': '2022-10-10',
    'dateTo': '2022-10-11',
    'aggregation': 'Actual'
}
 
meter_data_url = meter_data + timeseries_data['dateFrom'] + '/' + timeseries_data['dateTo'] + '/' + timeseries_data['aggregation']
 
meter_json = {
    "meteringPoints": {
        "meteringPoint": [
            first_meter
        ]
    }
}
 
meter_data_response = requests.post(meter_data_url, headers=headers, json=meter_json)

if meter_data_response.status_code != 200:
    print("Data lookup failed")
    print(meter_data_response.status_code)
    exit()

meter_json = json.loads(meter_data_response.text)
meter_json_formatted = json.dumps(meter_json, indent=2)
print(meter_json_formatted)

## TODO post to Humio
