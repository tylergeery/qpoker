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
        let type = event.target.getAttribute('data-type') || event.target.type;
        let value = event.target.value;

        for (let i = 0, l = keys.length; i < l; i++) {
            if ((i+1) == l) {
                state[keys[i]] = this.getValueType(type, value);
            }

            if (!state.hasOwnProperty(keys[i])) {
                state[keys[i]] = {};
            }

            state = state[keys[i]];
        }

        this.setState({values: this.state.values})
    }
}