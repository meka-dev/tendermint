# meka-dev/tendermint

This fork of [tendermint/tendermint](https://github.com/tendermint/tendermint)
has been patched to supoort the [Mekatek builder API](https://api.mekatek.xyz).
Each supported network has a tracking branch, which participating validators
should use when building their full node.

## Networks

<table style="text-align: left;">
  <tr>
    <th style="text-align: left;">Network</th>
    <th style="text-align: left;">Version</th>
    <th style="text-align: left;">Tendermint</th>
    <th style="text-align: left;">Fork</th>
    <th style="text-align: left;">Diff</th>
  </tr>
  <tr>
    <td><a href="https://github.com/osmosis-labs/osmosis">Osmosis</a></td>
    <td><a href="https://github.com/osmosis-labs/osmosis/tree/v11.0.1">v11.0.1</a></td>
    <td><a href="https://github.com/osmosis-labs/osmosis/blob/v11.0.1/go.mod#L28">v0.34.9</td>
    <td><a href="https://github.com/meka-dev/tendermint/tree/v0.0.0-osmosis-v11.0.1-a">v0.0.0-osmosis-v11.0.1-a</a></td>
    <td><a href="https://github.com/meka-dev/tendermint/compare/v0.34.9...v0.0.0-osmosis-v11.0.1-a">diff</a></td>
  </tr>
</table>

## Example

Here is an example of how to build Osmosis with a patched Tendermint.

```shell
git clone https://github.com/osmosis-labs/osmosis
cd osmosis
git checkout v11.0.1
go mod edit -replace=github.com/tendermint/tendermint=github.com/meka-dev/tendermint@v0.34.x-meka
go mod tidy
make install
```

________________ <br/> [upstream README](/README.upstream.md)
