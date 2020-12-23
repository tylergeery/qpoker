import * as React from "react";

type ErrorList = {
    errors: string[];
}

export class Errors extends React.Component<ErrorList, {}> {
    public render() {
        return this.props.errors.length ? (
            this.props.errors.map((err) => {
                return <div key={err}>
                    <p className="red-text text-lighten-1">{err}</p>
                </div>
            })
        ) : '';
    }
}