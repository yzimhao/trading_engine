// common/request.js

import CryptoJS from "crypto-js";
const appId = "1";
const baseUrl = "";

export const request = (url, data = {}, method = 'GET', parseBody=true) => {
  return new Promise((resolve, reject) => {
	let query = {};
	let nonce = Math.random().toString(36).slice(-8);
	
	query["t"] = Date.now();
	query["v"] = uni.getSystemInfoSync().version;
	
	
    let raw = ""
	if(method.toUpperCase()=="GET") {
        query = Object.assign({}, data, query);
        data = query;
        const sortedQuery = Object.keys(query).sort().map(key => `${key}=${query[key]}`).join('&');
		
        raw = method.toUpperCase() + "###" + url + "###" + sortedQuery + "###"  + "###" + nonce;
        data["sign"] = CryptoJS.MD5(raw).toString(CryptoJS.enc.Hex);
        // console.log("raw: ", raw, "sign: ", data["sign"]);
	}else if(method.toUpperCase() == "POST") {
        let postdatastr = JSON.stringify(data);
        const sortedQuery = Object.keys(query).sort().map(key => `${key}=${query[key]}`).join('&');
        raw = method.toUpperCase() + "###" + url + "###" + sortedQuery + "###" + postdatastr + "###" + nonce;
        query["sign"] = CryptoJS.MD5(raw).toString(CryptoJS.enc.Hex);
        // console.log("raw: ", raw, "sign: ", query["sign"]);
        url = url + "?" + Object.keys(query).sort().map(key => `${key}=${query[key]}`).join('&');
	}
					
    uni.request({
      url: baseUrl + url, // 拼接 URL
      method: method,
      data: data,
	  header:{
		  'X-N': nonce,
          'X-AppId': appId
	  },
      success: (res) => {
        if (res.statusCode === 200) {
          if(parseBody){
            if(res.data.code == 0) {
                resolve(res.data);
            }else{
                uni.showToast({
                    title: res.data.msg,
                    icon: "none"
                });
            }
          }
          if(!parseBody){
            resolve(res);
          }

        } else {
          if(res.data.code == 1000){
            uni.removeStorageSync('user');
          }
          reject(res);
        }
      },
      fail: (err) => {
        reject(err);
      }
    });
  });
};
