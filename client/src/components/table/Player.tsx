import * as React from "react";
import { Card, GamePlayer } from "../../objects/State";

type PlayerProps = {
    player: GamePlayer;
    index: number;
    gameState: string;
    cards: Card[];
}

type HandProps = {
    gameState: string;
    cards: Card[];
}

class Hand extends React.Component<HandProps, {}> {
    render() {
        return this.props.gameState != 'Init' ? (
            <div>
                <img className="card" src={`/assets/media/cards/${this.props.cards[0].imageName}.svg`} />
                <img className="card" src={`/assets/media/cards/${this.props.cards[0].imageName}.svg`} />
            </div>
        ) : '';
    }
}

export class Player extends React.Component<PlayerProps, {}> {
    render() {
        return <div className={ `player table-pos-${this.props.index}` }>
            {`Player ${this.props.player.id} ($${this.props.player.stack})` }
            <Hand gameState={this.props.gameState} cards={this.props.cards} />
        </div>
    }
}
