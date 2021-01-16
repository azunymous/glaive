import React from 'react';
import {Link} from "react-router-dom";
import objection from "./objection.gif";

function OpeningPost(props) {
    let {board, thread, imageContext} = props



    let displayImage = (post) => {
        let imgPath = post.thumbnail_image

        if (imageContext.startsWith("http")) {
            return imageContext + "/" + imgPath
        }
        return process.env.PUBLIC_URL + imageContext + imgPath;
    }


    let displayComment = (post, thread) => {
        if (post.comment_segments == null) {
            return post.comment;
        }

        return post.comment_segments.map((segment, i) => {
            return displaySegment(segment, i, post, thread)
        })
    }

    let displaySegment = (segment, i, post, thread) => {
        switch (segment.format[0]) {
            case "objection":
                return (
                    <div className={segment.format}><img src={objection} alt="Objection!"/></div>
                );
            case "roll":
                return (
                    <div className="roll">{Math.random() % 6}</div>
                );
            case "noQuote":
                return (
                    <div className={formatAsClasses(segment.format)} key={i}>
                        <span className="noQuote" key={i}>{segment.segment}</span>
                    </div>
                );
            default:
                return (
                    <div className={formatAsClasses(segment.format)} key={i}>{segment.segment}<br/></div>
                );
        }
    }

    let formatAsClasses = (segment) => {
        if (segment === null || segment.format === null) {
            return "";
        }
        return segment.map((format) => {
            return format + " "
        })
    }


    return <article className={"openingPost"}>
       <figure className={"catalogImage"}>
        <Link
            to={"/" + board + "/res/" + thread.post.no}>
            <span><img alt={thread.post.filename}
                                         src={displayImage(thread.post)}/></span>
        </Link>
       </figure>
        <figcaption className={"catalogComment"}>
            <span
                className="catalogSubject">{thread.subject}</span>
                {displayComment(thread.post, thread)}
        </figcaption>

    </article>


}



export default OpeningPost