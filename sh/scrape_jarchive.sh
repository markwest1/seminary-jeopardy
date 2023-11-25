#!/bin/bash

export data_dir=../data

# Usage: season_links [SEASON_ID_FILTER (as regular expression)]
season_links() {
    season_id_filter=$1
    slinks=$(grep ' href=\"showseason' "${data_dir}/listseasons.php" | sed 's/.*href=\"\([^\"]\+\)\".*/\1/')
    if [ -z $season_id_filter ]; then
        echo $slinks
    else
        for slink in $slinks
        do
            sid=$(echo $slink | sed 's/.*=\(.*\)/\1/')
            if [[ $sid =~ $season_id_filter ]]; then
                echo $slink
            fi
        done
    fi
}

# Usage: scrape_seasons SEASON_LINKS
scrape_seasons() {
    for slink in $1
    do
        sid=$(echo $slink | sed 's/.*=\(.*\)/\1/')
        if [[ $sid =~ ^[1-9]$ ]]; then
            sid=$(printf "%02d" $sid)
        fi
        echo curl --location "http://www.j-archive.com/$slink" -o "${data_dir}/seasons/season_${sid}.php"
    done
}

# Usage: scrape_games SEASON_PHP_FILES
scrape_games() {
    for season_file in $1
    do
        sid=$(echo $season_file | sed 's/.*season_\([^\.]\+\)\.php/\1/')
        game_folder="${data_dir}/seasons/season_${sid}"
        if [ ! -d $game_folder ]; then
            mkdir "$game_folder"
        fi

        echo "scraping_games from $season_file ..."
        game_hrefs=$(grep 'showgame' $season_file | sed 's/.*href=\"\([^\"]\+\)\".*/\1/')
        for href in $game_hrefs;
        do
            game_id=$(echo "$href" | sed 's/.*=//')
            game_php_file="${data_dir}/seasons/season_${sid}/game_${game_id}.php"
            echo "  scraping game $game_id to $game_php_file ..."
            curl -s --location $href -o $game_php_file
        done
        echo -------
    done
}

# export numerical_season_links=$(season_links ^[1-9][0-9]*$)
# scrape_seasons "$numerical_season_links"

# export alphabetic_season_links=$(season_links ^[A-Za-z_]\+$)
# scrape_seasons "$alphabetic_season_links"

export alphabetic_season_php_files="$(ls ${data_dir}/seasons/season_[a-z]*.php)"
scrape_games "$alphabetic_season_php_files"
