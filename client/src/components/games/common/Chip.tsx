import * as React from "react";
import { getChipAmount } from "../../../utils";

type ChipProps = {
    amount: number;
    color: "red" | "white"
};

export class Chip extends React.Component<ChipProps, {}> {
    render() {
        if (!this.props.amount) {
            return <div />;
        }

        return <div className={`chip ${this.props.color}`} title={this.props.amount.toString()}>
            {getChipAmount(this.props.amount)}
        </div>;
    }
}
