import * as React from "react";

import { classNames } from "../../utils";
import { Game } from "../../objects/Game";
import { ConnectionHandler } from "../../connection/ws";
import { EventState } from "../../objects/State";

type SettingsProps = {
    active: boolean;
    playerID: string;
    game?: Game;
    conn: ConnectionHandler;
    es: EventState;
}

type SettingsState = {
    options: any[];
    requests: any[];
    form: {
        [key: string]: any;
    }
}

export class Settings extends React.Component<SettingsProps, SettingsState> {
    constructor(props: any) {
        super(props)
        this.state = { options: [], requests: [], form: {
            'chip_request': 20000,
        }}
    }

    public componentDidMount() {
        this.props.conn.subscribe('admin', (event: any) => {
            console.log('settings received event:', event)
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
                    {isAdmin ? (
                        <tr>
                            <td colSpan={2}>Start Game:</td>
                            <td>
                                <button onClick={this.sendAdminAction.bind(this, 'start', '')} className="btn-flat green lighten-1" type="button">
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
                                <button onClick={this.sendChipsResponse.bind(this, req.id)} className="btn-flat green lighten-2" type="button">
                                    Approve
                                </button>
                            </td>
                            <td>
                                <button onClick={this.sendChipsResponse.bind(this, '-' + req.id)} className="btn-flat red lighten-2" type="button">
                                    Deny
                                </button>
                            </td>
                        </tr>
                    })}
                    <tr>
                        <td colSpan={2}>Request Chips</td>
                        <td>
                            <input type="number" step={50} defaultValue={20000} name="chip_request"
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
                    {this.state.options.map((option) => {
                        // TODO: render details if not admin
                        return <tr>
                            <td colSpan={2}>{option.label}:</td>
                            <td>
                                <input type={option.type} name={option.name} defaultValue={option.value}></input>
                            </td>
                        </tr>;
                    })}
                </table>
            </div>
        </div>;
    }
}