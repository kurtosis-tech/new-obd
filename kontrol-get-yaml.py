import requests
import yaml
import subprocess
import tempfile
import os


# Function to fetch JSON data from the endpoint
def fetch_json(endpoint_url):
    try:
        response = requests.get(endpoint_url)
        response.raise_for_status()  # Raise an exception for HTTP errors
        return response.json()
    except requests.RequestException as e:
        raise RuntimeError(f"HTTP request failed: {e}")


def fetch_yaml(endpoint_url):
    try:
        response = requests.get(endpoint_url)
        response.raise_for_status()  # Raise an exception for HTTP errors

        # Ensure the response content type is YAML
        if 'application/x-yaml' not in response.headers.get('Content-Type', ''):
            raise ValueError("Expected YAML response but got something else.")

        # Parse YAML response into Python objects (handling multiple documents)
        yaml_documents = list(yaml.safe_load_all(response.text))

        # Merge all YAML documents into a single YAML string
        merged_yaml = yaml.dump_all(yaml_documents, sort_keys=False)

        return merged_yaml
    except requests.RequestException as e:
        raise RuntimeError(f"HTTP request failed: {e}")
    except yaml.YAMLError as e:
        raise RuntimeError(f"Error parsing YAML: {e}")
    except ValueError as e:
        raise RuntimeError(f"Content type error: {e}")

# Function to convert JSON to YAML
def convert_json_to_yaml(json_data):
    try:
        # Clean the data to remove empty objects
        cleaned_data = remove_empty_objects(json_data)
        return yaml.dump(cleaned_data, sort_keys=False, default_flow_style=False)
    except yaml.YAMLError as e:
        raise RuntimeError(f"Error converting JSON to YAML: {e}")


# Function to remove empty JSON objects from the data
def remove_empty_objects(data):
    if isinstance(data, dict):
        return {k: remove_empty_objects(v) for k, v in data.items() if v not in ({}, [], None)}
    elif isinstance(data, list):
        return [remove_empty_objects(item) for item in data if item not in ({}, [], None)]
    else:
        return data


# Function to apply YAML data to the Kubernetes cluster
def apply_yaml_to_kubernetes(yaml_data):
    try:
        # Specify the full path to `kubectl` if not in PATH
        kubectl_path = 'kubectl'  # Adjust if necessary

        # Create a temporary file to store the YAML data
        with tempfile.NamedTemporaryFile(delete=False, mode='w', suffix='.yaml') as temp_file:
            temp_file.write(yaml_data)
            temp_file_path = temp_file.name

        # Run the `kubectl apply -f <file>` command
        subprocess.run([kubectl_path, 'apply', '-f', temp_file_path], check=True)

    finally:
        # Clean up the temporary file
        if os.path.exists(temp_file_path):
            os.remove(temp_file_path)


# Main function
def main():
    #endpoint_url = 'http://localhost:8080/tenant/7dca8131-d3fc-425e-af7e-2e0c94cc797f/cluster-resources'

    endpoint_url = 'http://localhost:8080/api/resources?uuid=7dca8131-d3fc-425e-af7e-2e0c94cc797f'

    try:

        # Fetch JSON data from the endpoint
        merged_yaml = fetch_yaml(endpoint_url)
        print(merged_yaml)  # Print the merged YAML document as a string

        # Extract the 'deployments' field from the JSON data
        # if 'deployments' in json_data:
        #     deployments = json_data['deployments']
        #     if len(deployments) > 3:
        #         deployments_data = deployments[3]
        #     else:
        #         raise ValueError("The 'deployments' list does not contain enough elements.")
        # else:
        #     raise ValueError("The 'deployments' field is not present in the JSON payload.")

        # Convert JSON data to YAML
        #yaml_data = convert_json_to_yaml(json_data)

        # Apply the YAML data to the Kubernetes cluster
        #apply_yaml_to_kubernetes(yaml_data)

        # Define the multi-line string with triple quotes
        #yaml_string = "apiVersion: apps/v1"


    except (RuntimeError, ValueError, subprocess.CalledProcessError) as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()
