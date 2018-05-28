import asyncio
from multiprocessing import Pool

import aiohttp


async def fetch(session, url):
    '''
    异步向中间件请求目标url地址
    '''
    async with session.get(url) as resp:
        if resp.status == 404:
            return None
        else:
            # 获取到的单个url
            return await response.text()


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
        await session.post(adress, data=bytes(data))

async def spider_html(adress):
    '''
    向目标地址获取html
    '''
    async with aiohttp.ClientSession() as session:
        async with session.get(adress) as res:
            print(await res.text())

async def task():
    '''
    各个任务的组装
    '''
    pass


def create_spider():
    '''
    爬虫程序的运行
    '''
    with Pool(processes=10) as pool:
        pass


if __name__ == '__main__':
    create_spider()
