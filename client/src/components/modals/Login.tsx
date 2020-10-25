import * as React from "react";
import * as Modal from "react-modal";

import { Errors } from "./Errors";

type LoginModalProps = {
    login: (values: object) => void;
    onRegisterClick: (event: any) => void;
}

type LoginModalState = {
    isOpen: boolean;
    errors: string[];
    values: {
        [key: string]: string;
    }
}

export class LoginModal extends React.Component<LoginModalProps, LoginModalState> {
    constructor(props: any) {
        super(props);

        this.state = { isOpen: false, errors: [], values: {}};
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
        event.preventDefault();
        this.props.login(this.state.values);
    }

    render() {
        return (
            <Modal
            isOpen={this.isActive()}
            onRequestClose={this.closeModal.bind(this)}
            contentLabel="Login Modal"
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

                <h2>Login</h2>

                <form onSubmit={this.submit.bind(this)}>
                    <Errors errors={this.state.errors} />
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
                        <button type="button" className="btn-large transparent grey-text text-darken-3 cancel-button" onClick={this.closeModal.bind(this)}>Close</button>
                        <button type="submit" className="btn-large grey darken-3">Submit</button>
                    </div>
                    <div className="row center">
                        <br/>
                        <a onClick={this.props.onRegisterClick} href="#signin">Sign up</a>
                    </div>
                </form>
            </Modal>
        );
    }
}