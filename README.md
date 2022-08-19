# meka-dev/tendermint

This fork of [tendermint/tendermint](https://github.com/tendermint/tendermint)
has been patched to supoort the [Mekatek builder API](https://api.mekatek.xyz).
Each supported network has a corresponding release tag, which participating
validators should use when building their full node.

<table>
  <tr>
    <th align="left">Network</th>
    <th align="left">Version</th>
    <th align="left">Tendermint</th>
    <th align="left">Tracking branch</th>
    <th align="left"><strong>Tagged release</strong></th>
    <th align="left">Diff</th>
  </tr>
  <tr>
    <td><a href="https://github.com/osmosis-labs/osmosis">Osmosis</a></td>
    <td><a href="https://github.com/osmosis-labs/osmosis/tree/v11.0.1">v11.0.1</a></td>
    <td><a href="https://github.com/osmosis-labs/osmosis/blob/v11.0.1/go.mod#L28">v0.34.9</td>
    <td><a href="https://github.com/meka-dev/tendermint/tree/v0.34.90-mekatek">v0.34.90-mekatek</a></td>
    <td><strong><a href="https://github.com/meka-dev/tendermint/tree/osmosis-v11.0.1-a">osmosis-v11.0.1-a</a><strong></td>
    <td><a href="https://github.com/meka-dev/tendermint/compare/v0.34.90...osmosis-v11.0.1-a">diff</a></td>
  </tr>
</table>

Here is an example of how to build Osmosis with a patched Tendermint.

```shell
git clone https://github.com/osmosis-labs/osmosis
cd osmosis
git checkout v11.0.1
go mod edit -replace=github.com/tendermint/tendermint=github.com/meka-dev/tendermint@osmosis-v11.0.1-a
go mod tidy
make install
```
