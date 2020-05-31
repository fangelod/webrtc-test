import axios from 'axios';

/**
 * Returns the data object from the response
 *
 * @param  {object} response A response from a network request
 *
 * @return {object}          The requested data
 */
function returnData(response) {
  return response.data;
}

/**
 * Checks if a network request came back fine, and throws an error if not
 *
 * @param  {object} response   A response from a network request
 *
 * @return {object|undefined} Returns either the response, or throws an error
 */
function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }

  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

/**
 * Requests a URL, returning a promise
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "fetch"
 *
 * @return {object}           The response data
 */
export default function request(url, options) {
  if (url[0] !== '/') {
    url = `api/${url}`;
  }

  options = options || {};
  options.credentials = 'same-origin';

  return axios(url, options).then(checkStatus).then(returnData);
}