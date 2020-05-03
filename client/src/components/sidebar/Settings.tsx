import * as React from "react";

import { classNames } from "../../utils";
import { Game } from "../../objects/Game";
import { ConnectionHandler } from "../../connection/ws";
import { EventState } from "../../objects/State";
import { UpdateGameRequest } from "../../requests/updateGame";
import { userStorage } from "../../utils/storage";

type SettingsProps = {
    active: boolean;
    playerID: string;
    game?: Game;
    conn: ConnectionHandler;
    es: EventState;
}

type SettingsState = {
    errors: string[];
    options: any[];
    requests: any[];
    form: {
        [key: string]: any;
    }
}

export class Settings extends React.Component<SettingsProps, SettingsState> {
    constructor(props: any) {
        super(props)
        this.state = {
            errors: [],
            options: Object.assign({}, props.game.options),
            requests: [],
            form: {
                'chip_request': props.game.options.buy_in_min
            }
        }
    }

    public componentDidMount() {
        this.props.conn.subscribe('admin', (event: any) => {
            if (!event.data.requests) {
                return
            }

            this.setState({requests: event.data.requests});
        });
    }

    public handleChange(event: React.ChangeEvent<HTMLInputElement>) {
        this.state.form[event.target.name] = event.target.value;
    }

    public sendChipsResponse(id: string) {
        this.sendAction('admin', { action: 'chip_response', value: id })
    }

    public sendAdminAction(action: string, valueKey: string) {
        let value = this.state.form[valueKey];

        this.sendAction('admin', { action, value })
    }

    public sendAction(type: string, data: any) {
        let action = {
            type,
            data
        }

        console.log("send action: ", action);
        this.props.conn.send(JSON.stringify(action))
    }

    public getOptions(): any[] {
        return [
            {
                name: 'big_blind',
                label: 'Big Blind',
                type: 'number',
                value: this.props.game.options.big_blind,
            },
            {
                name: 'capacity',
                label: 'Capacity',
                type: 'number',
                value: this.props.game.options.capacity,
            },
            {
                name: 'time_between_hands',
                label: 'Time Between Hands (s)',
                type: 'number',
                value: this.props.game.options.time_between_hands,
            },
            {
                name: 'buy_in_min',
                label: 'Min Buy In',
                type: 'number',
                value: this.props.game.options.buy_in_min,
            },
            {
                name: 'buy_in_max',
                label: 'Max Buy In',
                type: 'number',
                value: this.props.game.options.buy_in_max,
            }
        ]
    }

    private async updateGame(event: React.ChangeEvent<HTMLInputElement>) {
        switch (event.target.name) {
            case 'capacity':
                this.props.game.options.capacity = parseInt(event.target.value);
                break
            case 'big_blind':
                this.props.game.options.big_blind = parseInt(event.target.value);
                break
            case 'time_between_hands':
                this.props.game.options.time_between_hands = parseInt(event.target.value);
                break
            case 'buy_in_min':
                this.props.game.options.buy_in_min = parseInt(event.target.value);
                break
            case 'buy_in_max':
                this.props.game.options.buy_in_max = parseInt(event.target.value);
                break
        }

        const req = new UpdateGameRequest<Game>();
        const game = await req.request({
            id: this.props.game.id.toString(),
            data: this.props.game,
            userToken: userStorage.getToken(),
        });

        if (!req.success) {
            this.setState({errors: req.errors});
            return
        }
    }

    private isAdmin(): boolean {
        return this.props.game.owner_id.toString() == this.props.playerID;
    }

    public render() {
        let isAdmin = this.isAdmin();

        return <div className={classNames({"hidden": !this.props.active})}>
            <div>
                <table>
                    {isAdmin ? (
                        <tr>
                            <th colSpan={4}>Admin Control</th>
                        </tr>
                    ) : ''}
                    {(isAdmin && this.props.es.manager.status != "active") ? (
                        <tr>
                            <td colSpan={2}>Start Game:</td>
                            <td>
                                <button disabled={this.props.es.manager.status == "init"}
                                    onClick={this.sendAdminAction.bind(this, 'start', '')}
                                    className={classNames("btn-flat green lighten-1", {
                                        'disabled': this.props.es.manager.status == "init" 
                                    })} type="button">
                                    Start
                                </button>
                            </td>
                        </tr>
                    ) : ''}
                    {isAdmin ? (<tr></tr>) : ''}

                    <tr>
                        <th colSpan={4}>Chip Requests</th>
                    </tr>
                    {this.state.requests.map((req) => {
                        return <tr>
                            <td colSpan={2}>{this.props.es.getPlayer(req.player_id).username}</td>
                            <td>
                                <button onClick={this.sendChipsResponse.bind(this, req.player_id.toString())} className="btn-flat green lighten-2" type="button">
                                    Approve
                                </button>
                            </td>
                            <td>
                                <button onClick={this.sendChipsResponse.bind(this, '-' + req.player_id.toString())} className="btn-flat red lighten-2" type="button">
                                    Deny
                                </button>
                            </td>
                        </tr>
                    })}
                    <tr>
                        <td colSpan={2}>Request Chips</td>
                        <td>
                            <input type="number" step={this.props.game.options.big_blind} defaultValue={this.props.game.options.buy_in_min} name="chip_request"
                                onChange={this.handleChange.bind(this)} />
                        </td>
                        <td>
                            <button onClick={this.sendAdminAction.bind(this, 'chip_request', 'chip_request')} type="button" className="btn-flat">
                                Submit
                            </button>
                        </td>
                    </tr>
                   <tr></tr>

                    <tr></tr>
                    <th>Game Settings:</th>
                    {this.getOptions().map((option) => {
                        return isAdmin ? (<tr>
                            <td colSpan={2}>{option.label}:</td>
                            <td colSpan={2}>
                                <input onBlur={this.updateGame.bind(this)}
                                    type={option.type} name={option.name} defaultValue={option.value}></input>
                            </td>
                        </tr>) : (<tr>
                            <td colSpan={2}>{option.label}:</td>
                            <td colSpan={2}>{option.value}</td>
                        </tr>);
                    })}
                </table>
            </div>
        </div>;
    }
}