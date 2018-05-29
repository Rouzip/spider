import sys
import asyncio
from multiprocessing import Pool

import aiohttp
from requests_html import HTML

BASE_URL = 'https://api.ekyu.moe/d-spider/v1/'


class ExitException(Exception):
    '''
    退出的异常
    '''

    def __init__(self):
        super(ExitException, self).__init__()


async def fetch(session, url):
    '''
    异步向中间件请求目标url地址
    '''
    async with session.get(url) as resp:
        print(resp.status)
        if resp.status == 200 or resp.status == 204:
            # 获取到的单个url
            print(await resp.text(encoding='utf8'))
            return await resp.text(encoding='utf8')
        elif resp.status == 410:
            raise ExitException()
        else:
            raise Exception('other exception')


async def get_url(adress):
    '''
    向目标地址获取url
    '''
    async with aiohttp.ClientSession() as session:
        response = await fetch(session, adress)
        return response


async def post_html(adress, html):
    '''
    向目标地址发送已经读取到的html
    '''
    async with aiohttp.ClientSession() as session:
        await session.post(adress, data=bytes(html, encoding='utf8'))

async def spider_html(adress):
    '''
    向目标地址获取html
    '''
    async with aiohttp.ClientSession() as session:
        async with session.get(adress) as res:
            print(await res.text())
            return await res.text()

async def single_task():
    '''
    按照顺序各个请求任务
    '''
    try:
        url = await get_url(BASE_URL + 'URL')
        html = await spider_html(url)
        await post_html(BASE_URL + 'HTML', html)
    except Exception as e:
        return 'some exception'
    except ExitException as exit:
        # 再次抛出异常，让改loop退出
        raise ExitException
    finally:
        return 'lalala'


def callback(future):
    print('finish' + future.result())


def task():
    '''
    各个任务的组装
    '''
    while True:
        print('start')
        if sys.platform == 'win32':
            print('win')
            loop = asyncio.ProactorEventLoop()
            asyncio.set_event_loop(loop)
            task = asyncio.ensure_future(single_task())
            task.add_done_callback(callback)
            loop.run_until_complete(task)
            loop.close()
        else:
            print('linux')
            loop = asyncio.get_event_loop()
            task = asyncio.ensure_future(single_task())
            task.add_done_callback(callback)
            loop.run_until_complete(task)
            loop.close()


def create_spider():
    '''
    爬虫程序的运行
    '''
    with Pool(processes=10) as pool:
        pass


if __name__ == '__main__':
    task()
