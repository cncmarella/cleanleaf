/**
 * Single source of truth for company details that appear across the site and
 * in structured data. Update here, not in individual pages.
 */
export const site = {
  name: "CleanLeaf Technologies",
  shortName: "CleanLeaf",
  tagline: "Crop protection you can trust",
  description:
    "CleanLeaf Technologies manufactures and supplies crop protection products for farms across Telangana and beyond — formulated for effective pest control and responsible use.",
  url: "https://cleanleaf.vercel.app",
  address: {
    line1: "Sy.No: 1 & 3, Marripaliguda",
    line2: "Ghatkesar, Hyderabad - 501301",
    locality: "Hyderabad",
    region: "Telangana",
    postalCode: "501301",
    country: "IN",
  },
  phone: "8341099962",
  phoneHref: "tel:+918341099962",
  email: "cleanleaf789@gmail.com",
  emailHref: "mailto:cleanleaf789@gmail.com",
} as const;

export const nav = [
  { href: "/", label: "Home" },
  { href: "/products", label: "Products" },
  { href: "/about", label: "About" },
  { href: "/contact", label: "Contact" },
] as const;

/**
 * Placeholder catalogue. Replace with the real product range — names, actives,
 * pack sizes and registration numbers — once we have the official list.
 */
export const products = [
  {
    slug: "insecticides",
    name: "Insecticides",
    summary:
      "Contact and systemic formulations that bring sucking and chewing pests under control without setting the crop back.",
    crops: ["Cotton", "Paddy", "Chilli", "Vegetables"],
  },
  {
    slug: "fungicides",
    name: "Fungicides",
    summary:
      "Preventive and curative options for blight, mildew and leaf spot, built for the humid stretches of the season.",
    crops: ["Paddy", "Groundnut", "Grapes", "Tomato"],
  },
  {
    slug: "herbicides",
    name: "Herbicides",
    summary:
      "Selective and broad-spectrum weed control, timed for pre-emergence and post-emergence application windows.",
    crops: ["Paddy", "Maize", "Soybean", "Sugarcane"],
  },
  {
    slug: "plant-nutrition",
    name: "Plant Nutrition",
    summary:
      "Micronutrient mixes and growth promoters that help a treated crop recover quickly and fill out properly.",
    crops: ["All field crops", "Horticulture"],
  },
] as const;

export const strengths = [
  {
    title: "Made near your fields",
    body: "Our Ghatkesar facility keeps the supply chain short, so stock reaches dealers in days rather than weeks.",
  },
  {
    title: "Quality-checked batches",
    body: "Every batch is tested for active-ingredient concentration and stability before it leaves the plant.",
  },
  {
    title: "Guidance that fits the crop",
    body: "We advise on dosage, timing and tank mixing for the specific crop and pest, not a generic label instruction.",
  },
  {
    title: "Responsible use, always",
    body: "Clear labelling, safety guidance and pre-harvest intervals — so produce stays within residue limits.",
  },
] as const;
