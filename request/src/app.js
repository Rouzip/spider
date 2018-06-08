import axios from "axios";

const baseURL = "https://api.ekyu.moe/d-spider/v1/";
// const baseURL = 'http://localhost:3000/URL'

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
 * 默认维基不会出错
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
 * 默认中间件不会出错
 * @param {string} data
 */
async function post_data(data) {
  // 向中间件发送data
  try {
    console.log(data);
    let resp = await axios.post(baseURL + "HTML", data);
    if (resp.status === 200) throw new ExitError();
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
    let resp = await axios.get(adress);
    console.log(resp.data);
    if (resp.status === 410) throw new ExitError("停止爬虫");
    else if (resp.status === 200 || resp.status === 204) return resp.data;
  } catch (error) {
    console.log(error);
    throw new WaitError("wait a moment");
  }
}

/**
 * 将三个函数在这里组装
 */
async function task() {
  let url = await get_url();
  let resp = await request(url);
  await post_data(resp);
}

async function single_spider() {
  const promises = [];
  for (let i = 0; i < 3000; i++) promises.push(task());
  promises.forEach(task => {
    task.catch(err => {
      if (err instanceof ExitError) throw err;
    });
  });
}

/**
 * 尝试使用promise all来进行异步的组合
 */
async function spider_all() {
  try {
    for (let i = 0; i < 600; i++) await single_spider();
  } catch (error) {
    if (error instanceof ExitError) return "退出";
  }
}

async function task_single() {
  try {
    await task();
  } catch (error) {
    if (error instanceof ExitError) throw ExitError;
  }
}

async function task_500() {
  for (let i = 0; i < 500; i++) {
    await task_single();
  }
}

function main() {
  let index = 0;
  try {
    for (let i = 0; i < 800; i++) {
      console.log(index++);
      task_500()
        .then(res => {
          console.log(res);
        })
        .catch(err => {
          console.log(err);
          throw err;
        });
    }
  } catch (error) {
    if (error instanceof ExitError) {
      console.log("爬虫结束");
      return;
    }
  }
}

// main();
// task_500().catch(err => {
//   console.log(err);
// });
