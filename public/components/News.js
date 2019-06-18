import React from "react";
import moment from "moment";

class News extends React.Component {
    render() {
        var newsStyles = {
            float: "left",
            fontSize: "13px"
        };
        var newsRowStyles = {
            lineHeight: 1.4
        };

        if (this.props.news === undefined) {
            return null;
        }
        var news = this.props.news.items.map(item => (
            <div style={newsRowStyles}>
                {moment(item.pubDate).format("DD.MM.YYYY HH:mm")}{" "}
                <a href={item.rssLink} target="_blank">
                    {item.rssTitle.length > 100 ? item.rssTitle.slice(0, 100) + "..." : item.rssTitle}
                </a>
            </div>
        ));
        return (
            <div style={newsStyles}>
                <h2>Portfolio news</h2>
                {news}
            </div>
        );
    }
}

export default News;
