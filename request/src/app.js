import axios from "axios";

const baseURL = "https://api.ekyu.moe/d-spider/v1/";

/**
 * 自定义错误，退出爬虫
 */
class ExitError extends Error {
  constructor(message) {
    super(message);
    this.name = "ExitError";
  }
}

class WaitError extends Error {
  constructor(message) {
    super(message);
    this.name = "WaitError";
  }
}

/**
 * 获取目标网站的html数据，返回
 * @param {string} adress
 */
async function request(adress) {
  try {
    let resp = await axios.get(adress);
    if (resp.status === 200) return resp.data;
  } catch (error) {
    console.log(error);
  }
}

/**
 * 向中间件发送获取到的数据
 * @param {string} data
 */
async function post_data(data) {
  // 向中间件发送data
  try {
    let resp = await axios.post(baseURL + "HTML", data);
  } catch (error) {
    console.log(error);
  }
}

/**
 * 向服务器请求 需要爬取的url，每次获得的数据为一条url
 * 终止整个爬虫从这里终止，410终止
 */
async function get_url() {
  // 向中间件获取需要请求的URL
  try {
    let adress = baseURL + "URL";
    let resp = await axios.get(adress, { retry: 5, retryDelay: 1000 });
    if (resp.status === 410) throw new ExitError("停止爬虫");
    else if (resp.status === 200 || resp.status === 204) return resp.data;
  } catch (error) {
    throw new WaitError("wait a moment");
  }
}

/**
 * 将三个函数在这里组装
 */
async function task() {
  let index = 0;
  while (true) {
    try {
      console.log(index++);
      console.log(new Date());
      let url = await get_url();
      let resp = await request(url);
      await post_data(resp);
    } catch (error) {
      if (error instanceof ExitError) {
        return "退出爬虫";
      }
      console.log(error);
    }
  }
}

function main() {
  task()
    .then(resp => {
      console.log(resp);
    })
    .catch(error => {
      console.log(error);
    });
}

main();
