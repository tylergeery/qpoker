import * as React from "react";
import { Card, GamePlayer, Table } from "../../objects/State";
import { ConnectionHandler } from "../../connection/ws";
import { createGameAction } from "../../utils";

type PlayerProps = {
    playerID: string;
    player: GamePlayer;
    table: Table;
    index: number;
    gameState: string;
    cards: Card[];
    conn: ConnectionHandler;
}

type PlayerState = {
    timer: number;
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

export class Player extends React.Component<PlayerProps, PlayerState> {
    interval: number = null

    constructor(props: any) {
        super(props)
        this.interval = null;
        this.state = { timer: 0 };
    }

    private getCountdownTime(): number {
        let ts = +(new Date()) / 1000
        let seconds = 30 - (ts - this.props.table.activeAt);

        if (!seconds || seconds < 0) {
            seconds = 0;
        }

        return Math.floor(seconds);
    }
    private isSelected(): boolean {
        if (this.getCountdownTime() <= 0) {
            return false
        }

        return this.props.table.active === this.props.player.id;
    }

    private countDown(seconds: number) {
        if (seconds <= 0 || !this.isSelected()) {
            this.interval = null;
            this.setState({timer: 0})
            return;
        }

        this.setState({timer: seconds});
        window.setTimeout(this.countDown.bind(this, seconds-1), 1000);
    }

    private startTimer(): any {
        if (this.interval || !this.props.table.activeAt) {
            return
        }
        this.interval = 1;



        this.countDown(this.getCountdownTime());
    }

    render() {
        if (this.isSelected()) {
            this.startTimer()
        }

        return <div className={ `player table-pos-${this.props.index}` }>
            {`${this.props.player.username} (${this.props.player.stack})` }
            <Hand gameState={this.props.gameState} cards={this.props.cards} />
            <HandActions {...this.props} />
            <p>{this.state.timer ? this.state.timer.toString() : ''}</p>
        </div>
    }
}
