//go:build !noelk

package d2plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/util-go/xmain"
)

var ELKPlugin = elkPlugin{}

func init() {
	plugins = append(plugins, &ELKPlugin)
}

type elkPlugin struct {
	opts *d2elklayout.ConfigurableOpts
}

func (p elkPlugin) Flags(context.Context) ([]PluginSpecificFlag, error) {
	return []PluginSpecificFlag{
		{
			Name:    "elk-algorithm",
			Type:    "string",
			Default: d2elklayout.DefaultOpts.Algorithm,
			Usage:   "layout algorithm (e.g., layered)",
			Tag:     "elk.algorithm",
		},
		{
			Name:    "elk-edgeRouting",
			Type:    "string",
			Default: "",
			Usage:   "edge routing strategy (e.g., ORTHOGONAL, POLYLINE, SPLINES)",
			Tag:     "elk.edgeRouting",
		},
		{
			Name:    "elk-nodeNodeBetweenLayers",
			Type:    "int64",
			Default: int64(d2elklayout.DefaultOpts.NodeSpacing),
			Usage:   "spacing between nodes of adjacent layers",
			Tag:     "elk.layered.spacing.nodeNodeBetweenLayers",
		},
		{
			Name:    "elk-spacing-nodeNode",
			Type:    "int64",
			Default: int64(0),
			Usage:   "generic node-node spacing within a layer",
			Tag:     "elk.spacing.nodeNode",
		},
		{
			Name:    "elk-layered-spacing-nodeNodeBetweenLayers",
			Type:    "int64",
			Default: int64(d2elklayout.DefaultOpts.NodeSpacing),
			Usage:   "layered: spacing between nodes of adjacent layers",
			Tag:     "elk.layered.spacing.nodeNodeBetweenLayers",
		},
		{
			Name:    "elk-layered-spacing-edgeNodeBetweenLayers",
			Type:    "int64",
			Default: int64(d2elklayout.DefaultOpts.EdgeNodeSpacing),
			Usage:   "layered: spacing between edges and nodes of adjacent layers",
			Tag:     "elk.layered.spacing.edgeNodeBetweenLayers",
		},
		{
			Name:    "elk-spacing-edgeEdge",
			Type:    "int64",
			Default: int64(0),
			Usage:   "generic edge-edge spacing",
			Tag:     "elk.spacing.edgeEdge",
		},
		{
			Name:    "elk-layered-spacing-edgeEdgeBetweenLayers",
			Type:    "int64",
			Default: int64(0),
			Usage:   "layered: spacing between edges across adjacent layers",
			Tag:     "elk.layered.spacing.edgeEdgeBetweenLayers",
		},
		{
			Name:    "elk-padding",
			Type:    "string",
			Default: d2elklayout.DefaultOpts.Padding,
			Usage:   "padding for parent elements, e.g. [top=60,left=50,bottom=50,right=50]",
			Tag:     "elk.padding",
		},
		{
			Name:    "elk-layered-nodePlacement-strategy",
			Type:    "string",
			Default: "",
			Usage:   "layered node placement strategy (e.g., BRANDES_KOEPF)",
			Tag:     "elk.layered.nodePlacement.strategy",
		},
		{
			Name:    "elk-layered-nodePlacement-bk-fixedAlignment",
			Type:    "string",
			Default: "BALANCED",
			Usage:   "BK fixed alignment (e.g., BALANCED, LEFT, RIGHT, CENTER)",
			Tag:     "elk.layered.nodePlacement.bk.fixedAlignment",
		},
		{
			Name:    "elk-portConstraints",
			Type:    "string",
			Default: "",
			Usage:   "port constraints (e.g., FIXED_SIDE, FIXED_POS, FREE)",
			Tag:     "elk.portConstraints",
		},
		{
			Name:    "elk-libavoid-nudges",
			Type:    "bool",
			Default: false,
			Usage:   "enable libavoid nudging for edge routing",
			Tag:     "elk.edgeRouting.nudging",
		},
		{
			Name:    "elk-edgeLabel-spacing",
			Type:    "string",
			Default: "",
			Usage:   "edge label spacing (e.g., 10 or 10-14)",
			Tag:     "elk.edgeLabel.spacing",
		},
		{
			Name:    "elk-nodeLabels-placement",
			Type:    "string",
			Default: "",
			Usage:   "node labels placement (string per style guide)",
			Tag:     "elk.nodeLabels.placement",
		},
		{
			Name:    "elk-edgeNodeBetweenLayers",
			Type:    "int64",
			Default: int64(d2elklayout.DefaultOpts.EdgeNodeSpacing),
			Usage:   "spacing between nodes and edges routed next to the node’s layer",
			Tag:     "elk.layered.spacing.edgeNodeBetweenLayers",
		},
		{
			Name:    "elk-nodeSelfLoop",
			Type:    "int64",
			Default: int64(d2elklayout.DefaultOpts.SelfLoopSpacing),
			Usage:   "spacing between a node and its self loops",
			Tag:     "elk.spacing.nodeSelfLoop",
		},
	}, nil
}

func (p *elkPlugin) HydrateOpts(opts []byte) error {
	if opts != nil {
		var elkOpts d2elklayout.ConfigurableOpts
		err := json.Unmarshal(opts, &elkOpts)
		if err != nil {
			return xmain.UsageErrorf("non-ELK layout options given for ELK")
		}

		p.opts = &elkOpts
	}
	return nil
}

func (p elkPlugin) Info(ctx context.Context) (*PluginInfo, error) {
	opts := xmain.NewOpts(nil, nil)
	flags, err := p.Flags(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range flags {
		f.AddToOpts(opts)
	}
	return &PluginInfo{
		Name: "elk",
		Type: "bundled",
		Features: []PluginFeature{
			CONTAINER_DIMENSIONS,
			DESCENDANT_EDGES,
		},
		ShortHelp: "Eclipse Layout Kernel (ELK) with the Layered algorithm.",
		LongHelp: fmt.Sprintf(`ELK is a layout engine offered by Eclipse.
Originally written in Java, it has been ported to Javascript and cross-compiled into D2.
See https://d2lang.com/tour/elk for more.

Flags correspond to ones found at https://www.eclipse.org/elk/reference.html.

Flags:
%s
`, opts.Defaults()),
	}, nil
}

func (p elkPlugin) Layout(ctx context.Context, g *d2graph.Graph) error {
	return d2elklayout.Layout(ctx, g, p.opts)
}

func (p elkPlugin) PostProcess(ctx context.Context, in []byte) ([]byte, error) {
	return in, nil
}
