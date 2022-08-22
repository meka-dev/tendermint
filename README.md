# meka-dev/tendermint

This fork of [tendermint/tendermint](https://github.com/tendermint/tendermint)
has been patched to support the [Mekatek builder API](https://api.mekatek.xyz).
Each supported network has a tracking branch corresponding to their version of
Tendermint. Validators should build their nodes with this fork of Tendermint at
the corresponding release tag.

<table>
  <tr>
    <th>Network</th>
    <th>Version</th>
    <th>Tendermint</th>
    <th>Tracking branch</th>
    <th><strong> ★ Release tag ★ </strong></th>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/osmosis-labs/osmosis">Osmosis</a></td>
    <td align="center"><a href="https://github.com/osmosis-labs/osmosis/tree/v11.0.1">v11.0.1</a></td>
    <td align="center"><a href="https://github.com/osmosis-labs/osmosis/blob/v11.0.1/go.mod#L28">v0.34.19</td>
    <td align="center">
      <a href="https://github.com/meka-dev/tendermint/tree/v0.34.19-mekatek">v0.34.19-mekatek</a>
      (<a href="https://github.com/meka-dev/tendermint/compare/v0.34.19...v0.34.19-mekatek">diff</a>)
    </td>
    <td align="center">
      <strong><a href="https://github.com/meka-dev/tendermint/tree/osmosis-v11.0.1-a">osmosis-v11.0.1-a</a></strong>
      (<a href="https://github.com/meka-dev/tendermint/compare/v0.34.19...osmosis-v11.0.1-a">diff</a>)
    </td>
  </tr>
</table>

Here is an example of how to build Osmosis.

```shell
git clone https://github.com/osmosis-labs/osmosis
cd osmosis
git checkout v11.0.1
go mod edit -replace=github.com/tendermint/tendermint=github.com/meka-dev/tendermint@osmosis-v11.0.1-a
go mod tidy
make install
```
