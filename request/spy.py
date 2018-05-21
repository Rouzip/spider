import aiohttp

async def fetch(session, url):
    '''
    异步向中间件请求目标url地址
    '''
    async with session.get(url) as response:
        status = await response
        if status == 404:
            return None
        else:
            return await response.text()


async def get_url(adress):
    '''
    向目标地址获取url
    '''
    async with aiohttp.ClientSession() as session:
        response = await fetch(session, adress)
        return response


async def post_html(adress):
    '''
    向目标地址发送已经读取到的html
    '''
    async with aiohttp.ClientSession() as session:
        await session.get(adress)

def create_spider():
    '''
    主程序的运行
    '''
    pass

if __name__ == '__main__':
    pass
