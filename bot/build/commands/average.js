"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const discord_js_1 = require("discord.js");
const builders_1 = require("@discordjs/builders");
const api_1 = require("../api/api");
module.exports = {
    data: new builders_1.SlashCommandBuilder()
        .setName("average")
        .setDescription("Replies the average product for a price")
        .addStringOption((option) => option
        .setName("keyword")
        .setDescription("The product you want to search for")
        .setRequired(true)),
    async execute(interaction) {
        await interaction.deferReply();
        let userMessage = interaction.options.getString("keyword");
        let api = new api_1.Api();
        // @ts-ignore
        await api.init(process.env.key);
        let response = await getProductDetail(userMessage);
        let message = new discord_js_1.MessageEmbed()
            .setTitle(`Search term: ${userMessage}`)
            .setDescription(response)
            .setColor("RANDOM");
        await interaction.editReply({ embeds: [message] });
    },
    help: {
        name: "average",
        Description: "Retrieves the average price based on a search keyword",
        usage: "/average keyword: <search word>",
    },
};
let api;
async function init() {
    api = new api_1.Api();
    await api.init(process.env.key);
}
init();
async function getProductDetail(keyword) {
    let productData = await api.findNameByRegex(keyword);
    // console.log("getProductDetail productData: %s", JSON.stringify(productData)); // __AUTO_GENERATED_PRINT_VAR__
    let response = "";
    // if a single product is found
    if (productData.length == 1) {
        // get the first index of array
        let product = productData[0];
        response = `title: ${product.title}
      average price: ${product.average_price}
      average discount: ${product.average_discount}
      appearances: ${product.appearances}`;
    }
    // no product is found
    else if (productData.length == 0) {
        response = "No product found";
    }
    // if an array of product is found
    else {
        productData.forEach((element) => {
            let currentResponse = `title: ${element.title}
         average discount: ${element.average_discount}
         average price: ${element.average_price}
         appearances: ${element.appearances}

         `;
            response = response.concat(" ", currentResponse);
        });
    }
    return response;
}