#!/bin/bash

python manage.py collectstatic --noinput
python manage.py migrate
daphne bridge_constr.asgi:application --port $PORT --bind 0.0.0.0 -v2