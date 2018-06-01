import axios from "axios";

const baseURL = 'http://172.16.0.12:6001/'

/**
 * 自定义错误，退出爬虫
 */
class ExitError extends Error {
  constructor(message) {
    super(message)
    this.name = 'ExitError'
  }
}

/**
 * 获取目标网站的html数据，返回
 * @param {string} adress 
 */
async function request(adress) {
  try {
    let resp = await axios.get(adress)
    if (resp.status === 200) return resp.data
  } catch (error) {
    console.log(error)
    throw Error('')
  }
}

/**
 * 向中间件发送获取到的数据
 * @param {string} data 
 */
async function post_data(data) {
  // 向中间件发送data
  try {
    let resp = await axios.post(baseURL + 'POST', data)
  } catch (error) {
    console.log(error)
  }
}

/**
 * 向服务器请求 需要爬取的url，每次获得的数据为一条url
 * 终止整个爬虫从这里终止，410终止
 * FIXME: 暂停等待中间件使用403？
 */
async function get_url() {
  // 向中间件获取需要请求的URL
  try {
    let adress = baseURL + 'GET'
    let resp = await axios.get(adress)
      // FIXME: 添加状态，暂停向中间件请求url
    if (resp.status === 404)
      throw Error('暂停，等待url发送')
    else if (resp.data === 410)
      throw new ExitError('停止爬虫')
    return resp.data
  } catch (error) {
    console.log(error)
    if (error instanceof ExitError) {
      throw new ExitError()
    }
  }
  console.log(resp.data)
  return resp.data
}

/**
 * 将三个函数在这里组装
 */
async function task() {
  let index = 0
  while (true) {
    try {
      console.log(index++)
      let url = await get_url()
      let resp = await request(url)
      await post_data(resp)
    } catch (error) {
      return '退出爬虫'
    }
  }
}

function main() {
  task()
    .then(resp => {
      console.log(resp)
    })
    .catch(error => {
      console.log(error)
    })
}

main()







// while (true) {
console.log('start')
axios.get(baseURL + 'URL')
  .then(resp => {
    console.log(resp.data)
    let URL = resp.data
    axios.get(URL)
      .then(data => {
        axios.post(baseURL + 'HTML', data.data)
          .then(res => {
            console.log(res)
          })
      })
      .catch(e => {
        console.log(e)
      })
  })
  .catch(err => {
    console.log(err)
  })
  // }