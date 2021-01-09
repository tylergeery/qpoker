import * as React from "react";

import { classNames } from "../../utils";
import { ChipSettings } from "./ChipSettings";
import { Game, GameType } from "../../objects/Game";
import { ConnectionHandler, EventType } from "../../connection/ws";
import { UpdateGameRequest } from "../../requests/updateGame";
import { getGameType } from "../../utils/gameType";
import { userStorage } from "../../utils/storage";
import { AnonymousPlayer, AnonymousPlayerWithChips, findPlayer } from "../../objects/Player";

export type ManageButtonSettings = {
    text: string;
    disabled: boolean;
};

type SettingsProps = {
    active: boolean;
    game?: Game;
    players: AnonymousPlayer[]
    playerID: number;
    conn: ConnectionHandler;
    manageButtonSettings: ManageButtonSettings;
}

type SettingsState = {
    errors: string[];
    options: {
        [key: string]: any;
    };
    gameType?: GameType;
    requests: any[];
}

export class Settings extends React.Component<SettingsProps, SettingsState> {
    constructor(props: any) {
        super(props)
        this.state = {
            errors: [],
            options: Object.assign({}, props.game.options),
            gameType: null,
            requests: [],
        }
    }

    public componentDidMount() {
        this.props.conn.subscribe(EventType.admin, (event: any) => {
            if (!event.data.requests) {
                return
            }

            this.setState({requests: event.data.requests});
        });

        getGameType(this.props.game.game_type_id)
            .then((gameType: GameType) => {
                this.setState({gameType});
            });
    }

    public handleChange(event: React.ChangeEvent<HTMLInputElement>) {
        this.state.options[event.target.name] = event.target.value;
    }

    public clientAdminAction(action: string) {
        this.sendAction('admin', { action, value: {} })
    }

    public sendAction(type: string, data: any) {
        let action = {
            type,
            data
        }

        console.log("send action: ", action);
        this.props.conn.send(JSON.stringify(action))
    }

    private getOptions(): any[] {
        if (!this.state.gameType) {
            return [];
        }

        return this.state.gameType.options.map(opt => {
            return {
                name: opt.name,
                label: opt.label,
                type: opt.type,
                value: this.props.game.options[opt.name]
            };
        });
    }

    private getOption(name: string): any {
        let options = this.getOptions();

        for (let i=0; i < options.length; i++) {
            if (options[i].name === name) {
                return options[i];
            }
        }

        throw new Error("Unknown option name: " + name);
    }

    private async updateGame(event: React.ChangeEvent<HTMLInputElement>) {
        let option = this.getOption(event.target.name);

        switch(option.type) {
            case 'integer':
                this.props.game.options[option.name] = parseInt(event.target.value);
            case 'number':
                this.props.game.options[option.name] = parseFloat(event.target.value);
            default:
                this.props.game.options[option.name] = event.target.value;
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
        return this.props.game.owner_id.toString() == this.props.playerID.toString();
    }

    private supportsChips(): boolean {
        // TODO: do better
        return this.props.game.game_type_id === 1;
    }

    public render() {
        let isAdmin = this.isAdmin();
        let supportsChips = this.supportsChips();

        return <div className={classNames({"hidden": !this.props.active})}>
            <div>
                <table>
                    {isAdmin ? (
                        <tr>
                            <th colSpan={4}>Admin Control</th>
                        </tr>
                    ) : ''}
                    {(isAdmin) ? (
                        <tr>
                            <td colSpan={2}>{this.props.manageButtonSettings.text} Game:</td>
                            <td>
                                <button disabled={this.props.manageButtonSettings.disabled}
                                    onClick={this.clientAdminAction.bind(this, this.props.manageButtonSettings.text)}
                                    className={classNames("btn-flat green lighten-1", {
                                        'disabled': this.props.manageButtonSettings.disabled, 
                                    })} type="button">
                                    {this.props.manageButtonSettings.text}
                                </button>
                            </td>
                        </tr>
                    ) : ''}
                    {isAdmin ? (<tr></tr>) : ''}

                    {supportsChips ? (
                        <ChipSettings requests={this.state.requests}
                            player={findPlayer(this.props.playerID, this.props.players) as AnonymousPlayerWithChips}
                            players={this.props.players} game={this.props.game}
                            sendAction={this.sendAction.bind(this)} />
                    ) : ''}

                    <tr></tr>
                    <th>Game Settings:</th>
                    {this.getOptions().map((option) => {
                        return isAdmin ? (<tr key={option.name}>
                            <td colSpan={2}>{option.label}:</td>
                            <td colSpan={2}>
                                <input onChange={this.updateGame.bind(this)}
                                    type={option.type} name={option.name} defaultValue={option.value}></input>
                            </td>
                        </tr>) : (<tr key={option.name}>
                            <td colSpan={2}>{option.label}:</td>
                            <td colSpan={2}>{option.value}</td>
                        </tr>);
                    })}
                </table>
            </div>
        </div>;
    }
}
