import puppeteer from 'puppeteer';

const GOODREAD_URL = 'https://www.goodreads.com/book/show/';
const GOODREAD_IMAGE_URL_PATTERN = 'https://images-na.ssl-images-amazon.com/images';

export const find = (str, term, startsBy = 0) => {
  if (str === undefined) {
    return -1
  }
  let len = 0;
  let pos = null;
  for (let i = startsBy; i < str.length; i++) {
    if (str[i] == term[len])
      len++;
    else
      len = 0;
    if (len == term.length) {
      pos = i + 1 - term.length; //gets position i-term.length but has to add 1 given that startsBy has default value 0
      break;
    }
  }
  if (pos != null) return pos;

  return -1;
}

export const getGoodreadsUrl = (data) => {
  if (data === undefined) {
    return null;
  }
  let init = find(data, GOODREAD_URL);
  let final = find(data, "&", init + 10);
  let url = data.slice(init, final);
  return url;
}

export const getImageUrl = async (goodreadsLink) => {
  const browser = await puppeteer.launch({
    'args': [
      '--no-sandbox',
      '--disable-setuid-sandbox'
    ]
  });
  const page = await browser.newPage();
  await page.goto(goodreadsLink);
  const imageSelector = `img[src^="${GOODREAD_IMAGE_URL_PATTERN}"]`;
  await page.waitForSelector(imageSelector, {
    visible: true,
  });
  const src = await page.evaluate((selector) => {
    const img: any = document.querySelector(selector);
    if (!img) {
      throw new Error('Image not found');
    }
    return img.src;
  }, imageSelector);
  await browser.close();

  return src;
}