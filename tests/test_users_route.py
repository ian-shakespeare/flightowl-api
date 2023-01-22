import pytest
import sqlite3
import requests
import sys

sys.path.append('../')

USERS_URL = 'http://localhost:8000/users'
SESSIONS_URL = 'http://localhost:8000/sessions'

@pytest.fixture
def delete_test_user():
    yield
    conn = sqlite3.connect("../flightowl.db")
    cursor = conn.cursor()
    cursor.execute("DELETE FROM users WHERE email = 'test@email.com';")
    conn.commit()

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

def test_authenticate_user(create_test_user, delete_test_user):
    body = {'email': 'test@email.com',
            'password': 'password'}
    r = requests.post(SESSIONS_URL, json=body)
    assert r.status_code == 201
