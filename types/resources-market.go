package types

import (
	"errors"
	"math"
)

type MarketResourceID uint

const (
	OCTResourceID       MarketResourceID = 0
	MetalResourceID     MarketResourceID = 1
	CrystalResourceID   MarketResourceID = 2
	DeuteriumResourceID MarketResourceID = 3

	MAX_RESERVE_RATION float64 = 1
	N_RATIO            uint64  = 2
)

// based on https://www.linumlabs.com/articles/bonding-curves-the-what-why-and-shapes-behind-it
// price = ReserveRatio * resourceSupply ^ N_RATIO
// octBalance = resourceSupply^3 / 3
// priceForXToken = ((resourceSupply + X)^3 / 3) - (resourceSupply^3 / 3)

type OChainResourcesMarketWithAttributes struct {
	UniverseId            string  `cbor:"universeId"`
	FeesRate              float64 `cbor:"feesRate"`
	MetalReserveRatio     float64 `cbor:"metalReserveRatio"`
	MetalPoolBalance      uint64  `cbor:"metalPoolBalance"`
	MetalSupplyMinted     uint64  `cbor:"metalSupplyMinted"`
	CrystalReserveRatio   float64 `cbor:"crystalReserveRatio"`
	CrystalPoolBalance    uint64  `cbor:"crystalPoolBalance"`
	CrystalSupplyMinted   uint64  `cbor:"crystalSupplyMinted"`
	DeuteriumReserveRatio float64 `cbor:"deuteriumReserveRatio"`
	DeuteriumPoolBalance  uint64  `cbor:"deuteriumPoolBalance"`
	DeuteriumSupplyMinted uint64  `cbor:"deuteriumSupplyMinted"`
}

type OChainResourcesMarket struct {
	UniverseId string  `cbor:"1,keyasint"`
	FeesRate   float64 `cbor:"2,keyasint"`

	MetalReserveRatio float64 `cbor:"3,keyasint"`
	MetalPoolBalance  uint64  `cbor:"4,keyasint"`
	MetalSupplyMinted uint64  `cbor:"5,keyasint"`

	CrystalReserveRatio float64 `cbor:"6,keyasint"`
	CrystalPoolBalance  uint64  `cbor:"7,keyasint"`
	CrystalSupplyMinted uint64  `cbor:"8,keyasint"`

	DeuteriumReserveRatio float64 `cbor:"9,keyasint"`
	DeuteriumPoolBalance  uint64  `cbor:"10,keyasint"`
	DeuteriumSupplyMinted uint64  `cbor:"11,keyasint"`
}

func (market *OChainResourcesMarket) WithAttributes() OChainResourcesMarketWithAttributes {

	return OChainResourcesMarketWithAttributes{
		UniverseId:            market.UniverseId,
		FeesRate:              market.FeesRate,
		MetalReserveRatio:     market.MetalReserveRatio,
		MetalPoolBalance:      market.MetalPoolBalance,
		MetalSupplyMinted:     market.MetalSupplyMinted,
		CrystalReserveRatio:   market.CrystalReserveRatio,
		CrystalPoolBalance:    market.CrystalPoolBalance,
		CrystalSupplyMinted:   market.CrystalSupplyMinted,
		DeuteriumReserveRatio: market.DeuteriumReserveRatio,
		DeuteriumPoolBalance:  market.DeuteriumPoolBalance,
		DeuteriumSupplyMinted: market.DeuteriumSupplyMinted,
	}
}

func (market *OChainResourcesMarket) GetSwapAmountOut(from MarketResourceID, to MarketResourceID, amount uint64) (uint64, error) {

	//mint resource against woct
	if from == OCTResourceID {
		return market.calculateCurvedMintReturn(to, amount)
	}

	//burn resource for woct
	if to == OCTResourceID {
		return market.calculateCurvedBurnReturn(from, amount)
	}

	//burn resource aginst another resource
	burnReturn, err := market.calculateCurvedBurnReturn(from, amount)
	if err != nil {
		return 0, err
	}

	return market.calculateCurvedMintReturn(to, burnReturn)
}

func (market *OChainResourcesMarket) SwapResources(from MarketResourceID, to MarketResourceID, amount uint64) (uint64, error) {

	//mint resource against woct
	if from == OCTResourceID {
		return market.curvedMint(to, amount)
	}

	//burn resource for woct
	if to == OCTResourceID {
		return market.curvedBurn(from, amount)
	}

	//burn resource aginst another resource
	burnReturn, err := market.curvedBurn(from, amount)
	if err != nil {
		return 0, err
	}

	return market.curvedMint(to, burnReturn)
}

func (market *OChainResourcesMarket) calculateCurvedMintReturn(resourceId MarketResourceID, amount uint64) (uint64, error) {
	switch resourceId {
	case MetalResourceID:
		mintAmount, err := calculatePurchaseReturn(market.MetalSupplyMinted, market.MetalPoolBalance, market.MetalReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return mintAmount, nil
	case CrystalResourceID:
		mintAmount, err := calculatePurchaseReturn(market.CrystalSupplyMinted, market.CrystalPoolBalance, market.CrystalReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return mintAmount, nil
	case DeuteriumResourceID:
		mintAmount, err := calculatePurchaseReturn(market.DeuteriumSupplyMinted, market.DeuteriumPoolBalance, market.DeuteriumReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return mintAmount, nil
	default:
		return 0, errors.New("resourceIds dosn't exists")
	}
}

func (market *OChainResourcesMarket) calculateCurvedBurnReturn(resourceId MarketResourceID, amount uint64) (uint64, error) {

	switch resourceId {
	case MetalResourceID:
		returnedAmount, err := calculateSaleReturn(market.MetalSupplyMinted, market.MetalPoolBalance, market.MetalReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return returnedAmount, nil
	case CrystalResourceID:
		returnedAmount, err := calculateSaleReturn(market.CrystalSupplyMinted, market.CrystalPoolBalance, market.CrystalReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return returnedAmount, nil
	case DeuteriumResourceID:
		returnedAmount, err := calculateSaleReturn(market.DeuteriumSupplyMinted, market.DeuteriumPoolBalance, market.DeuteriumReserveRatio, amount)
		if err != nil {
			return 0, err
		}

		return returnedAmount, nil
	default:
		return 0, errors.New("resourceIds dosn't exists")
	}
}

func (market *OChainResourcesMarket) curvedMint(resourceId MarketResourceID, amount uint64) (uint64, error) {

	amountMinFees := uint64(float64(amount) * market.FeesRate)
	resourcesMinted, err := market.calculateCurvedMintReturn(resourceId, amount)
	if err != nil {
		return 0, err
	}

	switch resourceId {
	case MetalResourceID:

		market.MetalPoolBalance += amountMinFees
		market.MetalSupplyMinted += resourcesMinted

	case CrystalResourceID:
		market.CrystalPoolBalance += amountMinFees
		market.CrystalSupplyMinted += resourcesMinted

	case DeuteriumResourceID:
		market.DeuteriumPoolBalance += amountMinFees
		market.DeuteriumSupplyMinted += resourcesMinted

	default:
		return 0, errors.New("resourceIds dosn't exists")
	}

	return resourcesMinted, nil
}

func (market *OChainResourcesMarket) curvedBurn(resourceId MarketResourceID, amount uint64) (uint64, error) {

	octReturned, err := market.calculateCurvedBurnReturn(resourceId, amount)
	if err != nil {
		return 0, err
	}

	switch resourceId {
	case MetalResourceID:

		market.MetalPoolBalance -= octReturned
		market.MetalSupplyMinted -= amount

	case CrystalResourceID:
		market.CrystalPoolBalance -= octReturned
		market.CrystalSupplyMinted -= amount

	case DeuteriumResourceID:
		market.DeuteriumPoolBalance -= octReturned
		market.DeuteriumSupplyMinted -= amount

	default:
		return 0, errors.New("resourceIds dosn't exists")
	}

	amountMinFees := uint64(float64(octReturned) * market.FeesRate)
	return amountMinFees, nil
}

/*
*
  - @dev given a continuous token supply, reserve token balance, reserve ratio, and a deposit amount (in the reserve token),
  - calculates the return for a given conversion (in the continuous token)
    *
  - @param supply           continuous token total supply
  - @param reserveBalance   total reserve token balance
  - @param reserveRatio     reserve ratio, represented in ppm, 1-1000000
  - @param depositAmount    deposit amount, in reserve token
    *
  - @return purchase return amount: supply * ((1 + _depositAmount / _reserveBalance) ^ (_reserveRatio / MAX_RESERVE_RATIO) - 1)
    PurchaseReturn = ContinuousTokenSupply * ((1 + ReserveTokensReceived / ReserveTokenBalance) ^ (ReserveRatio - 1)
*/
func calculatePurchaseReturn(supply uint64, reserveBalance uint64, reserveRatio float64, amount uint64) (uint64, error) {
	if reserveRatio == 1 {
		return supply * amount / reserveBalance, nil
	}

	return uint64(
		float64(supply) * math.Pow(
			(1+float64(amount))/float64(reserveBalance),
			reserveRatio-1,
		),
	), nil
}

/**
 * @dev given a continuous token supply, reserve token balance, reserve ratio and a sell amount (in the continuous token),
 * calculates the return for a given conversion (in the reserve token)
 *
 * Formula:
 * Return = reserveBalance * (1 - (1 - amount / supply) ^ (1 / reserveRatio))
 *
 * @param _supply              continuous token total supply
 * @param _reserveBalance    total reserve token balance
 * @param _reserveRatio     constant reserve ratio, represented in ppm, 1-1000000
 * @param _sellAmount          sell amount, in the continuous token itself
 *
 * @return sale return amount
 */
func calculateSaleReturn(supply uint64, reserveBalance uint64, reserveRatio float64, amount uint64) (uint64, error) {
	return uint64(
		float64(reserveBalance) * (1 - math.Pow(1.0-float64(amount)/float64(supply), 1/reserveRatio)),
	), nil
}

// Continuous Token Price = Reserve Token Balance / (Continuous Token Supply x Reserve Ratio)
func ContiniousTokenPrice(supply uint64, reserveBalance uint64, reserveRatio uint64) float64 {
	return float64(reserveBalance) / (float64(supply) * float64(reserveRatio))
}
