const GOODREAD_URL = "https://www.goodreads.com/book/show/";
const GOODREAD_IMAGE_URL_PATTERN =
  "https://images-na.ssl-images-amazon.com/images";

export const find = (htmlResponse: string, pattern: string, startsBy = 0) => {
  if (htmlResponse === undefined) {
    return -1;
  }

  let len = 0;
  let pos = null;
  for (let i = startsBy; i < htmlResponse.length; i++) {
    if (htmlResponse[i] == pattern[len]) len++;
    else len = 0;
    if (len == pattern.length) {
      pos = i + 1 - pattern.length; //gets position i-pattern.length but has to add 1 given that startsBy has default value 0
      break;
    }
  }
  if (pos != null) return pos;

  return -1;
};

export const getGoodreadsUrl = (response: string) => {
  if (response === undefined) {
    return null;
  }

  let init = find(response, GOODREAD_URL);
  let final = find(response, "&", init + 10);
  let goodreadsUrl = response.slice(init, final);
  return goodreadsUrl;
};

export const getImageUrl = (response: string) => {
  if (response === undefined) {
    return null;
  }

  let init = find(response, GOODREAD_IMAGE_URL_PATTERN);
  let final = find(response, '"', init + 10);
  let imageUrl = response.slice(init, final);
  return imageUrl;
};
