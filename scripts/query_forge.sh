#!/bin/bash


URL_FILES="https://forgeapi.puppetlabs.com"
URL_RELEASES="https://forgeapi.puppetlabs.com/v3/releases"
EXCLUDE_FIELDS="exclude_fields=readme,changelog,license,reference,tags,malware_scan,metadata"
INCREMENT=10


download_module_release()
{
    local module="$1"

    curl -s -O "$URL_FILES/$module"
}


get_number_of_modules()
{
    local query="$1"
    local limit=1

    total=$(curl -s "${URL_RELEASES}?limit=${limit}&${query}&${EXCLUDE_FIELDS}" \
        | jq ".pagination.total")

    if [ -z "$total" ] || [ "$total" == "null" ];
    then
        echo "0"
    else
        echo "$total"
    fi
}


# list_module_release(query) lists all modules for the <query>
# Examples
#    list_module_releases "owner=puppetlabs"
#    list_module_releases "module=puppetlabs-stdlib"
#
list_module_releases()
{
    local query="$1"

    number_of_modules=$(get_number_of_modules "$query")

    local first=0
    local increment="$INCREMENT"
    local last="$number_of_modules"

    for current in $(seq "$first" "$increment" "$last")
    do
        curl -s "${URL_RELEASES}?offset=${current}&limit=${increment}&${query}&${EXCLUDE_FIELDS}" \
            | jq ".results | .[].file_uri" \
            | sed 's/"//g'
    done
}


in_cache()
{
    local filename="$1"
    if [ -f "$filename" ];
    then
	    return 1
    fi
    return 0
}


download_modules_from()
{
    local filename="$1"

    while IFS= read -r file_uri
    do
        module=$(basename "$file_uri")

        if in_cache "$module";
        then
            echo -n "downloading $file_uri ... "
            download_module_release "$file_uri"
	        echo "done"
        else
            echo "already in cache $file_uri"
        fi
    done < "$filename"
}


usage()
{
    echo "Usage: $0 [-l <query> ] [-d <file_uri>] [ -D <filename> ]" 1>&2
    exit 1
}


while getopts ":l:d:D:" options; do
    case "${options}" in
        l)
            query=${OPTARG}
            #echo "list $query"
            list_module_releases "$query"
            ;;
        d)
            file_uri=${OPTARG}
            #echo "download $file_uri"
            download_module_release "$file_uri"
            ;;
        D)
            filename=${OPTARG}
            #echo "download from files $filename"
            download_modules_from "$filename"
            ;;

        *)
            usage
            ;;
    esac
done
