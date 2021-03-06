# purpose of this file: scrape penguin magic open box section for coin products
# Date: 2021-09-03
# ---------------------------------
import requests
from bs4 import BeautifulSoup
import os
import sys
import json
# from pprint import pprint

cwd = os.getcwd()
parent_dir = os.path.abspath(os.path.join(cwd, os.pardir))
# print(f" parent_dir: {str(parent_dir)}")  # __AUTO_GENERATED_PRINT_VAR__

sys.path.insert(
    1, parent_dir
)  # add parent directory to system path in this program to access utils
import utils


def get_webpage(product: dict):
    """

    Args:
        product (dict): product dictiry

    Returns:
        BeautifulSoup object from the webpage
    """
    url = product["url"]
    # get web html
    html_page = requests.get(url).text
    return BeautifulSoup(html_page, "html.parser")


def validate(product: dict):
    """
    makes sure the product is a coin product

    Args:
        product (dict): the product to be validated
    """
    # NOTE: condition should be `and`
    if ("coin" or "coins") not in product["description"].lower() and (
        "coin" or "coins"
    ) not in product["title"].lower():
        print(product["title"], "is not a coin product")
        exit(1)

    with open(
        os.path.join(os.path.dirname(__file__), "product_info.txt"), "+r"
    ) as file:
        if file.read() != product["title"]:
            #  print("product changed")

            # overwrite current product title
            with open(
                os.path.join(os.path.dirname(__file__), "product_info.txt"), "w"
            ) as file:
                file.write(product["title"])
        else:
            # product is the same
            raise AttributeError("Product has not changed")


def get_product_info(soup, product):
    """
    get the product info of a product

    Args:
        soup (): The soup object of the product
        product (): The product dictionary

    Returns:
        None
    """
    try:
        # find product title
        product["title"] = soup.find("div", id="product_name").h1.text  # type:ignore
        # product description
        product["description"] = soup.find("div", class_="product_subsection").text.format()  # type: ignore
    except:
        raise AttributeError("There are no open box products currently")
        #  print("There are no open box product currently")


def main():
    product = {"title": str, "description": str, "url": str}

    # product["url"] = "https://www.penguinmagic.com/openbox/"
    product["url"] = "https://www.penguinmagic.com/p/1797"

    # create soup object
    soup = get_webpage(product=product)

    # get product title and description
    get_product_info(soup, product)

    if not utils.if_interested(product["title"]):
        print("Product is not interesting...")
        exit(1)

    # current product is a coin product
    validate(product)

    # delete the description
    del product["description"]

    # print product dictionary to stdout
    print(json.dumps(product))


if __name__ == "__main__":
    main()
