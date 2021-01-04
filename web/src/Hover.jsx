import React from "react";


export function Hover({onHover, children}) {
    return (
        <span className="hover">
                <span className="hover__no-hover">{children}</span>
                <span className="hover__hover">{onHover}</span>
            </span>
    );
}
