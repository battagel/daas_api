#!/usr/bin/env python3
from time import gmtime, strftime
import json

###################################################################################
# This file should process data.json into a DaaS compatible format. This should
# then be outputted into a new file which can be used later when calling the API
# with a bulk_import
#
# Data structure:
# [ // Multiple words stored
#     {
#         "phrase": "", // This is later used as the key.
#         "terms": [""],  // Any search terms to include. This can/should include code variables.
#         "last_update": "", // For term up-to-date-ness
#         "complexity": 0.0, // Float. How basic or complex the term. Good for refining search later
#         "topics": [""], // Index for related topics / projects
#         "explanations": [ // Multiple potential explanations
#             {
#                 "definition": "", // Any HTMl text here
#                 "author": "", // Last author of the explanation
#                 "tags": [""], // List of tags. This can include code languages, projects, components.
#                 "code": [""], // List of explanatory code snippets.
#                 "references": [""], // Multiple references. What about files?
#                 "heat": 0.0 // Float. Public vote for how accurate / relevant / descriptive the explanation is. List in accuracy order.
#             }
#         ]
#     }
# ]
###################################################################################

input_filepath = "data.json"
output_filepath = "processed_data.json"

def read_file():
    with open(input_filepath, "r") as file:
        data = json.load(file)
        return data


def convert_data(data):
    # Strip of metadata
    new_data = []
    for item in data:
        new_item = {
            "phrase": item["word_key"],
            "terms": ["HPE"],
            "last_update": strftime("%Y-%m-%d %H:%M:%S", gmtime()),
            "complexity": 0.0,
            "topics": [],
            "explanations": []
        }
        for i, definition in enumerate(item["definitions"]):
            explanation = {
                "definition": definition,
                "author": "Matthew Battagel",
                "tags": [],
                "code": [],
                "references": item["references"][i],
                "heat": 0.0
            }
            new_item["explanations"].append(explanation)
        new_data.append(new_item)
    return new_data


def write_file(data):
    with open(output_filepath, "w", encoding="utf-8") as file:
        json.dump(data, file, ensure_ascii=False, indent=4)

def run():
    print("## Starting ##")
    data = read_file()
    new_data = convert_data(data["data"]["data"][0]["rows"])
    write_file(new_data)
    print("## Success! ##")


if __name__ == "__main__":
    run()
