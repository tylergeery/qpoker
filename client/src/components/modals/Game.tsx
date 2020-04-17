import * as React from "react";
import * as Modal from "react-modal";

import { BaseModal } from "./Base";
import { Errors } from "./Errors";
import { Game } from "../../objects/Game";
import { CreateGameRequest } from "../../requests/createGame";
import { userStorage } from "../../utils/storage";

type GameModalProps = {
    game?: Game;
}

export class GameModal extends BaseModal<GameModalProps> {
    constructor(props: any) {
        super(props)

        this.state.values.options = { capacity: 9, big_blind: 50 };
    }

    async submit(event: any) {
        const req = new CreateGameRequest<Game>();
        const game = await req.request({
            data: this.state.values,
            userToken: userStorage.getToken(),
        })
    
        if (!req.success) {
            this.setState({errors: req.errors});
            return
        }

        window.location.href = `/${game.slug}`;

        return
    }

    render() {
        return (
            <Modal
            isOpen={this.isActive()}
            onRequestClose={this.closeModal.bind(this)}
            contentLabel="Game Modal"
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

                <h2>{this.props.game && this.props.game.id ? 'Update' : 'Create'} Game</h2>

                <form>
                    <Errors errors={this.state.errors} />
                    <div>
                        <label>
                            Game Name:
                            <input onChange={this.recordValue.bind(this)} type="text" name="name" />
                        </label>
                    </div>
                    <div>
                        <label>
                            Capacity:
                            <input onChange={this.recordValue.bind(this)} type="number" defaultValue="9" step="1" name="options.capacity" />
                        </label>
                    </div>
                    <div>
                        <label>
                            Big Blind:
                            <input onChange={this.recordValue.bind(this)} type="number" defaultValue="50" step="2" name="options.big_blind" />
                        </label>
                    </div>
                    <div className="row center">
                        <br/>
                        <button type="button" className="btn-large transparent grey-text text-darken-3" onClick={this.closeModal.bind(this)}>Close</button>
                        <button type="button" className="btn-large grey darken-3" onClick={this.submit.bind(this)}>Submit</button>
                    </div>
                </form>
            </Modal>
        );
    }
}