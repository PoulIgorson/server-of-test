from django.shortcuts import render
frm django.http import JsonResponse

from datetime import date


def index_page(request):
    return render(request, 'index.html', {'pagename': 'Главная'})


def today_json(request):
    return JsonResponse({'today': f'{date.today()}'})
