import * as React from "react";

interface ModalState {
    isOpen: boolean;
    errors: string[];
    ctx: {
        [key: string]: any;
    };
    values: {
        [key: string]: any;
    }
}

export abstract class BaseModal<P> extends React.Component<P, ModalState> {
    constructor(props: any) {
        super(props)

        this.state = {isOpen: false, errors: [], ctx: {}, values: {}};
    }

    openModal() {
        this.setState({isOpen: true});
    }

    closeModal() {
        this.setState({isOpen: false, errors: []});
    }

    isActive() {
        return this.state && this.state.isOpen;
    }

    getValueType(type: string, value: string): any {
        switch(type) {
            case "number":
                return parseInt(value, 10);
            default:
                return value;
        }
    }

    recordValue(event: any) {
        let keys = event.target.name.split('.');
        let state = this.state.values;

        for (let i = 0, l = keys.length; i < l; i++) {
            if ((i+1) == l) {
                state[keys[i]] = this.getValueType(event.target.type, event.target.value);
            }

            if (!state.hasOwnProperty(keys[i])) {
                state[keys[i]] = {};
            }

            state = state[keys[i]];
        }
    }
}