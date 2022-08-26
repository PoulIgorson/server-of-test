from django.shortcuts import render

from datetime import date


def index_page(request):
    return render(
        request,
        'index.html', {
            'pagename': 'Главная',
            'today': f'{date.today()}'
        }
    )
