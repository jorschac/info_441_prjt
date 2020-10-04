import React, { Component } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus} from '@fortawesome/free-solid-svg-icons';

export class CardContainer extends Component {
    render() {
        let paletteCards = this.props.propList.filteredData.map(x => <PaletteCard widget={x} handleSelectPalette={this.props.propList.handleClick} 
            handleResetLock={this.props.propList.handleResetLock}/>)

        return (
            <section id="cardcontainer">
                {paletteCards}
            </section>
        );
    }
}


class AddButton extends Component {
    /*trackInput = (e) => {
        this.props.propList.handleSearch(e.target.value);
    }*/

    handleClick = (e) => {
        /*e.preventDefault();
        let filter = this.props.propList.searchQuery;

        if (this.props.propList.filterList.includes(filter)) {
            this.props.propList.handleError('Already filtered!');
            
        } else {
            let selectedColorNames = this.props.propList.selectedPalette.map(x => convert.hex.keyword(x));
            let lockId = selectedColorNames.indexOf(filter);
            if (lockId !== -1) {
                this.props.propList.handleLock(filter, lockId);
            }
            this.props.propList.handleAddFilter(filter);
        }*/
    }

    render() {
        return (
                <button type="submit" id="addbutton" onClick={this.handleClick}>
                    <p>add a new wedget</p>
                    <FontAwesomeIcon icon={faPlus} className='fas fa-plus' aria-hidden="true" />
                </button>
            
        );
    }
}

class PaletteCard extends Component {
    handleClick = () => {
        /*let colors = [this.props.palette.light_shade, this.props.palette.light_accent, this.props.palette.main, 
            this.props.palette.dark_accent, this.props.palette.dark_shade];*/
        this.props.handleResetLock();
        //this.props.handleSelectPalette(colors);
    }

    render() {
        /*let colors = [this.props.palette.light_shade, this.props.palette.light_accent, this.props.palette.main, 
            this.props.palette.dark_accent, this.props.palette.dark_shade];*/
        
        //colors = colors.map(x => <PaletteCardColor color={x} />);
        //console.log(this.props.widget.WidgetCoverPic);
        return (
            <div className="palette" aria-label="color palette" onClick={this.handleClick}  style={{backgroundImage: `url(${this.props.widget.WidgetCoverPic})`}}>
                <div className="setinfo">
                    <div className="colorcontainer">
                    <p className="WigetTitle">{this.props.widget.WigetTitle}</p>
                    <p className="WigetDesc">{this.props.widget.WigetDesc}</p>                                      
                    </div>
                </div>
            </div>
        );
    }
}



/*
class PaletteCardColor extends Component {
    render() {
        let bgcolor = {backgroundColor: this.props.color};
        return (
            <div className="color" style={bgcolor}></div>
        );
    }
}*/
/*
export class NumberOfResult extends Component {
    render() {
        return (
            <p id="nPalettes" aria-label="number of search results" aria-live="polite">
                {this.props.nResult + ' results found'}
            </p>
        );
    }
}
*/