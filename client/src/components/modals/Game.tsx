import * as React from "react";
import * as Modal from "react-modal";

import { BaseModal } from "./Base";
import { Errors } from "./Errors";
import { Game, GameType, Option } from "../../objects/Game";
import { CreateGameRequest } from "../../requests/createGame";
import { userStorage } from "../../utils/storage";
import { getGameTypes } from "../../utils/gameType";

type GameModalProps = {
    game?: Game;
}

export class GameModal extends BaseModal<GameModalProps> {
    constructor(props: any) {
        super(props)

        this.state.ctx.types = [];
        this.state.values.options = {};
    }

    public componentDidMount() {
        this.fetchGameTypes()
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

    private fetchGameTypes() {
        if (this.state.ctx.types.length) {
            return
        }

        if (!userStorage.getID()) {
            return
        }

        getGameTypes()
            .then(resp => {
                if (resp) {
                    this.setState({ctx: {types: resp}})
                }
            }, err => {
                console.error('error fetching game types: ', err)
            });
    }

    private getSelectedIndex() {
        let gameTypeID = this.state.values.game_type_id;

        for (let i=0; i < this.state.ctx.types.length; i++) {
            if (this.state.ctx.types[i].id == gameTypeID) {
                return i;
            }
        }

        return -1;
    }

    private getOptions(): Option[] {
        let selected = this.getSelectedIndex()

        if (selected === -1) {
            return []
        }

        return this.state.ctx.types[selected].options;
    }

    private getOptionType(opt: Option): string {
        switch(opt.type) {
            case 'number':
            case 'integer':
                return 'number';
            default:
                return 'text';
        }
    }

    render() {
        this.fetchGameTypes();

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
                            Game Type:
                            <select className="show-select" onChange={this.recordValue.bind(this)} name="game_type_id" data-type="number">
                                <option>Select Game Type</option>
                                {this.state.ctx.types.map((type: GameType) =>
                                    <option key={type.id} value={type.id}>
                                        {type.display_name}
                                    </option>
                                )}
                            </select>
                        </label>
                    </div>
                    <div>
                        <label>
                            Game Name:
                            <input onChange={this.recordValue.bind(this)} type="text" name="name" />
                        </label>
                    </div>
                    {this.getOptions().map((opt: Option) => {
                        return <div key={`${opt.game_type_id}-${opt.game_option_id}`}>
                            <label>
                                {opt.label}:
                                <input onChange={this.recordValue.bind(this)} type={this.getOptionType(opt)}
                                        defaultValue={opt.default_value} name={`options.${opt.name}`}/>
                            </label>
                        </div>;
                    })}
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