import * as React from "react";
import { Card, GamePlayer } from "../../objects/State";
import { ConnectionHandler } from "../../connection/ws";
import { createGameAction } from "../../utils";

type PlayerProps = {
    playerID: string;
    player: GamePlayer;
    index: number;
    gameState: string;
    cards: Card[];
    conn: ConnectionHandler;
}

type HandProps = {
    gameState: string;
    cards: Card[];
}

class Hand extends React.Component<HandProps, {}> {
    render() {
        return this.props.gameState != 'Init' ? (
            <div>
                <img className="card" src={`/assets/media/cards/${this.props.cards[0].imageName()}.svg`} />
                <img className="card" src={`/assets/media/cards/${this.props.cards[1].imageName()}.svg`} />
            </div>
        ) : '';
    }
}

class HandActions extends React.Component<PlayerProps, {}> {
    bet: number;

    constructor(props: any) {
        super(props)
        this.bet = 0;
    }

    public getOptions(): string[] {
        let options: string[] = [];

        for (let opt in this.props.player.options) {
            if (this.props.player.options[opt]) {
                options.push(opt.slice(4));  // Remove `can_` prefix
            }
        }

        return options;
    }

    public sendAction(event: React.ChangeEvent<HTMLInputElement>) {
        switch (event.target.innerHTML) {
            case 'bet':
                this.props.conn.send(createGameAction({
                    name: event.target.innerHTML,
                    amount: this.bet,
                }));
                break;
            default:
                this.props.conn.send(createGameAction({
                    name: event.target.innerHTML,
                    amount: 0,
                }));
                break;
        }
    }

    private registerBet(event: React.ChangeEvent<HTMLInputElement>) {
        this.bet = parseInt(event.target.value);
    }

    render() {
        let options = this.getOptions();

        return this.props.playerID.toString() === this.props.player.id.toString() ? <div>
            {options.map((opt) => {
                return <button onClick={this.sendAction.bind(this)} type="button">{opt}</button>;
            })}
            {options.length ? <input type="number" defaultValue={0} onChange={this.registerBet.bind(this)} /> : ''}
        </div> : ''
    }
}

export class Player extends React.Component<PlayerProps, {}> {
    render() {
        return <div className={ `player table-pos-${this.props.index}` }>
            {`${this.props.player.username} (${this.props.player.stack})` }
            <Hand gameState={this.props.gameState} cards={this.props.cards} />
            <HandActions {...this.props} />
        </div>
    }
}
