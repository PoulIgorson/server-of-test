from django.shortcuts import render

from datetime import date


def index_page(request):
    return render(request, 'index.html', {'pagename': 'Главная'})


def today_json(request):
    return {
        'status': 200,
        'body': {
            'today': f'{date.today()}'
        }
    }
