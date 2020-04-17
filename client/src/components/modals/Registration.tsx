import * as React from "react";
import * as Modal from "react-modal";

import { Errors } from "./Errors";

type RegisterModalProps = {
    register: (values: object) => void;
    onLoginClick: (event: any) => void;
}

type RegisterModalState = {
    isOpen: boolean;
    errors: string[];
    values: {
        [key: string]: string;
    }
}

export class RegistrationModal extends React.Component<RegisterModalProps, RegisterModalState> {
    constructor(props: any) {
        super(props)

        this.state = {isOpen: false, errors: [],  values: {}};
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

    recordValue(event: any) {
        this.state.values[event.target.name] = event.target.value;
    }

    submit(event: any) {
        this.props.register(this.state.values);
    }

    render() {
        return (
            <Modal
            isOpen={this.isActive()}
            onRequestClose={this.closeModal.bind(this)}
            contentLabel="Registration Modal"
            style={{
                overlay: {
                    backgroundColor: 'rgba(0, 0, 0, 0.7)'
                },
                content: {
                    margin: 'auto',
                    maxHeight: '600px',
                    minWidth: '350px',
                    maxWidth: '600px',
                    width: '50%',
                    color: 'lightsteelblue'
                }
            }}
            >

                <h2>Register</h2>

                <form>
                    <Errors errors={this.state.errors} />
                    <div>
                        <label>
                            Username:
                            <input onChange={this.recordValue.bind(this)} type="text" name="username" />
                        </label>
                    </div>
                    <div>
                        <label>
                            Email:
                            <input onChange={this.recordValue.bind(this)} type="text" name="email" />
                        </label>
                    </div>
                    <div>
                        <label>
                            Password:
                            <input onChange={this.recordValue.bind(this)} type="password" name="pw" />
                        </label>
                    </div>
                    <div className="row center">
                        <br/>
                        <button type="button" className="btn-large transparent grey-text text-darken-3" onClick={this.closeModal.bind(this)}>Close</button>
                        <button type="button" className="btn-large grey darken-3" onClick={this.submit.bind(this)}>Submit</button>
                    </div>
                    <div className="row center">
                        <br/>
                        <a onClick={this.props.onLoginClick} href="#login">Login</a>
                    </div>
                </form>
            </Modal>
        );
    }
}