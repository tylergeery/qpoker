import * as React from "react";

type SeatProps = {
    index: number;
};

export class Seat extends React.Component<SeatProps, {}> {
    render() {
        return <div className={ `seat table-pos-${this.props.index}` }></div>
    }
}
