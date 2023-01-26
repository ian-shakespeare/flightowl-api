import pytest
import requests
import os
from datetime import datetime

USERS_URL = 'http://localhost:8000/users'
SESSIONS_URL = 'http://localhost:8000/sessions'
FLIGHTS_URL = 'http://localhost:8000/flights'
TESTING_URL = 'http://localhost:8000/tests'
TESTING_KEY = os.environ["TESTING_KEY"]

@pytest.fixture
def delete_test_user():
    yield
    body = {'token': TESTING_KEY}
    return requests.delete(TESTING_URL, json=body)

@pytest.fixture
def create_test_user():
    body = {'firstName': 'First',
            'lastName': 'Last',
            'email': 'test@email.com',
            'password': 'password',
            'sex': 'male'}
    return requests.post(USERS_URL, json=body)

def test_create_user(create_test_user, delete_test_user):
    assert create_test_user.status_code == 201

def test_create_existing_user(create_test_user, delete_test_user):
    create_test_user.status_code == 409

def test_authenticate_user(create_test_user, delete_test_user):
    body = {'email': 'test@email.com',
            'password': 'password'}
    r = requests.post(SESSIONS_URL, json=body)
    assert r.status_code == 201

def test_get_flights(create_test_user, delete_test_user):
    cookies = create_test_user.cookies
    body = {'originLocationCode': 'LAS',
            'destinationLocationCode': 'LAX',
            'departureDate': '2023-06-06',
            'adults': 1}
    res = requests.get(FLIGHTS_URL, json=body, cookies=cookies)
    assert res.status_code == 200
