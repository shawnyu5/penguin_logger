import { Client, TextChannel } from "discord.js";
// const exec = require("child_process").execSync;
import { execSync } from "child_process";
import { ICoinProduct } from "./types/coinProduct";
import { environment } from "./enviroments/enviroment";

/**
 * return a channel information by name
 * @param client - discord client
 * @param channelName - name of channel to search for
 * @returns channel information with the name passed in. If not found. undefined
 */
function getChannelByName(
   client: Client,
   channelName: string
): TextChannel | undefined {
   const channel = client.channels.cache.find((ch) => {
      // @ts-ignore
      return ch.name == channelName;
   });
   return channel as TextChannel;
}

/**
 * check if the current product is a coin product
 * @returns return json parsed string from `coin_products.py`. Other wise return null
 */
function checkCoinProduct(): ICoinProduct | null {
   // let result = execSync("python3 ../coin_products/coin_products.py");
   try {
      let result = execSync(
         "python3 ../coin_products/coin_products.py"
      ).toString();
      // console.log("checkCoinProduct result.toString(): %s", result); // __AUTO_GENERATED_PRINT_VAR__
      result = result.split("{")[1];
      result = "{" + result;
      // console.log(JSON.parse(result)); // __AUTO_GENERATED_PRINT_VAR__
      return JSON.parse(result);
   } catch (error) {
      // console.log(error);
      return null;
   }
}

/**
 * generates a message pinging all users in config.json about the coinProduct
 * @param coinProduct - the coin product
 * @returns A message string
 */
function buildMessage(coinProduct: ICoinProduct): string {
   let message: string = "COIN PRODUCT ALERT! ";
   message += `<@${environment.COIN_PRODUCT_ALERT_USERS}> `;
   message += `

title: ${coinProduct.Title}
url: https://www.penguinmagic.com/openbox/`;
   return message;
}

export { checkCoinProduct, getChannelByName, buildMessage };
