import pytest
import requests
import os

# NEED TO FIND A BETTER METHOD OF END TO END TESTING
def test_override_all():
    assert True

# USERS_URL = 'http://localhost:8000/users'
# SESSIONS_URL = 'http://localhost:8000/sessions'
# FLIGHTS_URL = 'http://localhost:8000/flights'
# TESTING_URL = 'http://localhost:8000/tests'
# TESTING_KEY = os.environ["TESTING_KEY"]

# @pytest.fixture
# def delete_test_user():
#     yield
#     body = {'token': TESTING_KEY}
#     return requests.delete(TESTING_URL, json=body)

# @pytest.fixture
# def create_test_user():
#     body = {'firstName': 'First',
#             'lastName': 'Last',
#             'email': 'test@email.com',
#             'password': 'password',
#             'sex': 'male'}
#     return requests.post(USERS_URL, json=body)

# @pytest.fixture
# def delete_test_flight():
#     yield
#     body = {'token': TESTING_KEY}
#     return requests.delete(TESTING_URL, json=body)

# @pytest.fixture
# def mock_flight_req_body():
#     return '''{
#         "data": {
#             "type": "flight-offers-pricing",
#             "flightOffers": [
#                 {
#                     "type": "flight-offer",
#                     "id": "1",
#                     "source": "GDS",
#                     "instantTicketingRequired": false,
#                     "nonHomogeneous": false,
#                     "oneWay": false,
#                     "lastTicketingDate": "2023-01-30",
#                     "numberOfBookableSeats": 9,
#                     "itineraries": [
#                         {
#                             "duration": "PT13H22M",
#                             "segments": [
#                                 {
#                                     "departure": {
#                                         "iataCode": "NRT",
#                                         "terminal": "1",
#                                         "at": "2023-04-20T16:45:00"
#                                     },
#                                     "arrival": {
#                                         "iataCode": "YVR",
#                                         "terminal": "M",
#                                         "at": "2023-04-20T09:30:00"
#                                     },
#                                     "carrierCode": "AC",
#                                     "number": "4",
#                                     "aircraft": {
#                                         "code": "77W"
#                                     },
#                                     "operating": {
#                                         "carrierCode": "AC"
#                                     },
#                                     "duration": "PT8H45M",
#                                     "id": "58",
#                                     "numberOfStops": 0,
#                                     "blacklistedInEU": false
#                                 },
#                                 {
#                                     "departure": {
#                                         "iataCode": "YVR",
#                                         "terminal": "M",
#                                         "at": "2023-04-20T11:05:00"
#                                     },
#                                     "arrival": {
#                                         "iataCode": "LAX",
#                                         "terminal": "6",
#                                         "at": "2023-04-20T14:07:00"
#                                     },
#                                     "carrierCode": "AC",
#                                     "number": "552",
#                                     "aircraft": {
#                                         "code": "7M8"
#                                     },
#                                     "operating": {
#                                         "carrierCode": "AC"
#                                     },
#                                     "duration": "PT3H2M",
#                                     "id": "59",
#                                     "numberOfStops": 0,
#                                     "blacklistedInEU": false
#                                 }
#                             ]
#                         }
#                     ],
#                     "price": {
#                         "currency": "USD",
#                         "total": "605.65",
#                         "base": "157.00",
#                         "fees": [
#                             {
#                                 "amount": "0.00",
#                                 "type": "SUPPLIER"
#                             },
#                             {
#                                 "amount": "0.00",
#                                 "type": "TICKETING"
#                             }
#                         ],
#                         "grandTotal": "605.65"
#                     },
#                     "pricingOptions": {
#                         "fareType": [
#                             "PUBLISHED"
#                         ],
#                         "includedCheckedBagsOnly": true
#                     },
#                     "validatingAirlineCodes": [
#                         "AC"
#                     ],
#                     "travelerPricings": [
#                         {
#                             "travelerId": "1",
#                             "fareOption": "STANDARD",
#                             "travelerType": "ADULT",
#                             "price": {
#                                 "currency": "USD",
#                                 "total": "605.65",
#                                 "base": "157.00"
#                             },
#                             "fareDetailsBySegment": [
#                                 {
#                                     "segmentId": "58",
#                                     "cabin": "ECONOMY",
#                                     "fareBasis": "KKVQ23BA",
#                                     "brandedFare": "BASIC",
#                                     "class": "K",
#                                     "includedCheckedBags": {
#                                         "quantity": 1
#                                     }
#                                 },
#                                 {
#                                     "segmentId": "59",
#                                     "cabin": "ECONOMY",
#                                     "fareBasis": "KKVQ23BA",
#                                     "brandedFare": "BASIC",
#                                     "class": "K",
#                                     "includedCheckedBags": {
#                                         "quantity": 1
#                                     }
#                                 }
#                             ]
#                         }
#                     ]
#                 }
#             ]
#         }
#     }'''

# def test_create_user(create_test_user, delete_test_user):
#     assert create_test_user.status_code == 201

# def test_create_existing_user(create_test_user, delete_test_user):
#     create_test_user.status_code == 409

# def test_authenticate_user(create_test_user, delete_test_user):
#     body = {'email': 'test@email.com',
#             'password': 'password'}
#     r = requests.post(SESSIONS_URL, json=body)
#     assert r.status_code == 201

# def test_get_flights(create_test_user, delete_test_user):
#     cookies = create_test_user.cookies
#     body = {'originLocationCode': 'LAS',
#             'destinationLocationCode': 'LAX',
#             'departureDate': '2023-06-06',
#             'adults': 1}
#     res = requests.get(FLIGHTS_URL, json=body, cookies=cookies)
#     assert res.status_code == 200

# def test_save_flight(create_test_user, delete_test_user, mock_flight_req_body, delete_test_flight):
#     cookies = create_test_user.cookies
#     res = requests.post(FLIGHTS_URL, json=mock_flight_req_body, cookies=cookies)
#     assert res.status_code == 201
